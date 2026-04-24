package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/term"
)

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func connectToAD(hostname string, port int, username, password string, skipVerify bool) (*ldap.Conn, error) {
	// TLS configuration
	tlsConfig := &tls.Config{
		ServerName:         hostname,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: skipVerify,
	}

	conn, err := ldap.DialURL(
		fmt.Sprintf("ldaps://%s:%d", hostname, port),
		ldap.DialWithTLSConfig(tlsConfig),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial AD server: %w", err)
	}

	if err := conn.Bind(username, password); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("bind failed: %w", err)
	}

	return conn, nil
}

func readPassword(bindUser string) (string, error) {
	fmt.Printf("LDAP Password for %s: ", bindUser)
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func searchUserBySAMAccountName(conn *ldap.Conn, baseDN, samAccountName string) (*UserResult, error) {
	// Escape user input to prevent LDAP injection
	escapedSAM := ldap.EscapeFilter(samAccountName)

	req := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(sAMAccountName=%s)", escapedSAM),
		[]string{
			"dn",
			"cn",
			"sAMAccountName",
			"displayName",
			"mail",
		},
		nil,
	)

	res, err := conn.Search(req)
	if err != nil {
		return nil, fmt.Errorf("LDAP search failed: %w", err)
	}

	if len(res.Entries) == 0 {
		return nil, ldap.NewError(
			ldap.LDAPResultNoSuchObject,
			fmt.Errorf("no user found with sAMAccountName=%s", samAccountName),
		)
	}

	e := res.Entries[0]
	return &UserResult{
		DN:             e.DN,
		CN:             e.GetAttributeValue("cn"),
		SAMAccountName: e.GetAttributeValue("sAMAccountName"),
		DisplayName:    e.GetAttributeValue("displayName"),
		Mail:           e.GetAttributeValue("mail"),
	}, nil
}

func generateConfigFile(path string, cfg Config) error {

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() { _ = file.Close() }()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func main() {
	// CLI flags
	configPath := flag.String("config", "", "Path to JSON config file")
	bindUserFlag := flag.String("b", "", "AD bind user (DOMAIN\\\\user)")
	sam := flag.String("sam", "", "sAMAccountName to search for")
	jsonOut := flag.Bool("json", false, "Output result as JSON")
	jsonGen := flag.Bool("json-gen", false, "Generate example JSON config file")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification (use with caution)")
	flag.Parse()

	// Defaults
	cfg := Config{
		Hostname: "hostname",
		Port:     636,
		BaseDN:   "DC=ad,DC=company,DC=com",
	}

	if *jsonGen {
		err := generateConfigFile("example.json", cfg)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("✅ example.json generated successfully")
		os.Exit(0)
	}

	if *sam == "" {
		log.Fatal("Missing required flag: -sam")
	}

	// Load config file if provided
	if *configPath != "" {
		fileCfg, err := loadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
		cfg = *fileCfg
	}

	// Override config with CLI flags
	if *bindUserFlag != "" {
		cfg.BindUser = *bindUserFlag
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Password handling
	password := os.Getenv("LDAP_PASSWORD")
	if password == "" {
		var err error
		password, err = readPassword(cfg.BindUser)
		if err != nil {
			log.Fatal("Failed to read password:", err)
		}
	}

	if *insecure {
		log.Println("⚠️  WARNING: TLS certificate verification is disabled")
	}

	conn, err := connectToAD(cfg.Hostname, cfg.Port, cfg.BindUser, password, *insecure)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	user, err := searchUserBySAMAccountName(conn, cfg.BaseDN, *sam)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok &&
			ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			if *jsonOut {
				fmt.Println(`{"error":"user not found"}`)
			} else {
				fmt.Println("❌ No user found")
			}
			return // Exit without os.Exit to allow defer to run
		}
		// Print error and return to allow defer to run
		log.Println(err)
		return
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(user); err != nil {
			log.Printf("Failed to encode JSON output: %v", err)
			return
		}
	} else {
		fmt.Println("✅ User found")
		fmt.Printf("DN: %s\n", user.DN)
		fmt.Printf("CN: %s\n", user.CN)
		if user.SAMAccountName != "" {
			fmt.Printf("sAMAccountName: %s\n", user.SAMAccountName)
		}
		if user.DisplayName != "" {
			fmt.Printf("DisplayName: %s\n", user.DisplayName)
		}
		if user.Mail != "" {
			fmt.Printf("Email: %s\n", user.Mail)
		}
	}
}
