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

type Config struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	BaseDN   string `json:"baseDN"`
	BindUser string `json:"bindUser"`
}

type UserResult struct {
	DN             string `json:"dn"`
	CN             string `json:"cn"`
	SAMAccountName string `json:"sAMAccountName"`
	DisplayName    string `json:"displayName,omitempty"`
	Mail           string `json:"mail,omitempty"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func connectToAD(hostname string, port int, username string, password string) (*ldap.Conn, error) {
	addr := fmt.Sprintf("%s:%d", hostname, port)

	conn, err := ldap.DialTLS("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dial AD server: %w", err)
	}

	if err := conn.Bind(username, password); err != nil {
		conn.Close()
		return nil, fmt.Errorf("bind failed: %w", err)
	}

	return conn, nil
}

func readPassword() (string, error) {
	fmt.Print("LDAP Password: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func searchUserBySAMAccountName(conn *ldap.Conn, baseDN string, samAccountName string) (*UserResult, error) {
	req := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(sAMAccountName=%s)", samAccountName),
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

func main() {
	// CLI flags
	configPath := flag.String("config", "", "Path to JSON config file")
	bindUserFlag := flag.String("b", "", "AD bind user (DOMAIN\\\\user)")
	sam := flag.String("sam", "", "sAMAccountName to search for")
	jsonOut := flag.Bool("json", false, "Output result as JSON")
	flag.Parse()

	if *sam == "" {
		log.Fatal("Missing required flag: -sam")
	}

	// Defaults
	cfg := Config{
		Hostname: "localhost",
		Port:     636,
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

	if cfg.Hostname == "" || cfg.Port == 0 || cfg.BaseDN == "" || cfg.BindUser == "" {
		log.Fatal("Incomplete configuration (hostname, port, baseDN, bindUser required)")
	}

	// Password handling
	password := os.Getenv("LDAP_PASSWORD")
	if password == "" {
		var err error
		password, err = readPassword()
		if err != nil {
			log.Fatal("Failed to read password:", err)
		}
	}

	conn, err := connectToAD(cfg.Hostname, cfg.Port, cfg.BindUser, password)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	user, err := searchUserBySAMAccountName(conn, cfg.BaseDN, *sam)
	if err != nil {
		if ldapErr, ok := err.(*ldap.Error); ok &&
			ldapErr.ResultCode == ldap.LDAPResultNoSuchObject {
			if *jsonOut {
				fmt.Println(`{"error":"user not found"}`)
			} else {
				fmt.Println("❌ No user found")
			}
			os.Exit(1)
		}
		log.Fatal(err)
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(user)
	} else {
		fmt.Println("✅ User found")
		fmt.Printf("DN: %s\n", user.DN)
		fmt.Printf("CN: %s\n", user.CN)
		fmt.Printf("sAMAccountName: %s\n", user.SAMAccountName)
		fmt.Printf("Display Name: %s\n", user.DisplayName)
		fmt.Printf("Email: %s\n", user.Mail)
	}
}
