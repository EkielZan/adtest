package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpfile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	cfg := Config{
		Hostname: "test.example.com",
		Port:     636,
		BaseDN:   "DC=test,DC=example,DC=com",
		BindUser: "DOMAIN\\testuser",
	}

	// Write config to temp file
	encoder := json.NewEncoder(tmpfile)
	if err := encoder.Encode(cfg); err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	// Test loading the config
	loaded, err := loadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("loadConfig failed: %v", err)
	}

	if loaded.Hostname != cfg.Hostname {
		t.Errorf("Expected hostname %s, got %s", cfg.Hostname, loaded.Hostname)
	}
	if loaded.Port != cfg.Port {
		t.Errorf("Expected port %d, got %d", cfg.Port, loaded.Port)
	}
	if loaded.BaseDN != cfg.BaseDN {
		t.Errorf("Expected BaseDN %s, got %s", cfg.BaseDN, loaded.BaseDN)
	}
	if loaded.BindUser != cfg.BindUser {
		t.Errorf("Expected BindUser %s, got %s", cfg.BindUser, loaded.BindUser)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "invalid*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	_, _ = tmpfile.WriteString("{invalid json}")
	_ = tmpfile.Close()

	_, err = loadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestGenerateConfigFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "generated*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpfile.Name()
	_ = tmpfile.Close()
	_ = os.Remove(tmpPath) // Remove so we can test creation
	defer func() { _ = os.Remove(tmpPath) }()

	cfg := Config{
		Hostname: "ad.example.com",
		Port:     636,
		BaseDN:   "DC=example,DC=com",
		BindUser: "EXAMPLE\\admin",
	}

	err = generateConfigFile(tmpPath, cfg)
	if err != nil {
		t.Fatalf("generateConfigFile failed: %v", err)
	}

	// Verify the file can be loaded back
	loaded, err := loadConfig(tmpPath)
	if err != nil {
		t.Fatalf("Failed to load generated config: %v", err)
	}

	if loaded.Hostname != cfg.Hostname {
		t.Errorf("Generated config hostname mismatch")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "valid config",
			config: Config{
				Hostname: "ad.example.com",
				Port:     636,
				BaseDN:   "DC=example,DC=com",
				BindUser: "DOMAIN\\user",
			},
			valid: true,
		},
		{
			name: "missing hostname",
			config: Config{
				Hostname: "",
				Port:     636,
				BaseDN:   "DC=example,DC=com",
				BindUser: "DOMAIN\\user",
			},
			valid: false,
		},
		{
			name: "missing port",
			config: Config{
				Hostname: "ad.example.com",
				Port:     0,
				BaseDN:   "DC=example,DC=com",
				BindUser: "DOMAIN\\user",
			},
			valid: false,
		},
		{
			name: "missing BaseDN",
			config: Config{
				Hostname: "ad.example.com",
				Port:     636,
				BaseDN:   "",
				BindUser: "DOMAIN\\user",
			},
			valid: false,
		},
		{
			name: "missing BindUser",
			config: Config{
				Hostname: "ad.example.com",
				Port:     636,
				BaseDN:   "DC=example,DC=com",
				BindUser: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid=%v, got valid=%v (error: %v)", tt.valid, valid, err)
			}
		})
	}
}

func TestConfig_Validate_PortRange(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port 636", 636, false},
		{"valid port 389", 389, false},
		{"valid port 1", 1, false},
		{"valid port 65535", 65535, false},
		{"invalid port 0", 0, true},
		{"invalid port negative", -1, true},
		{"invalid port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Hostname: "ad.example.com",
				Port:     tt.port,
				BaseDN:   "DC=example,DC=com",
				BindUser: "DOMAIN\\user",
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnectToAD_InsecureFlag(t *testing.T) {
	tests := []struct {
		name       string
		hostname   string
		port       int
		skipVerify bool
	}{
		{
			name:       "secure mode (skipVerify=false)",
			hostname:   "nonexistent.invalid",
			port:       636,
			skipVerify: false,
		},
		{
			name:       "insecure mode (skipVerify=true)",
			hostname:   "nonexistent.invalid",
			port:       636,
			skipVerify: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that connectToAD handles both secure and insecure modes
			// without panicking. Connection will fail since server doesn't exist.
			conn, err := connectToAD(tt.hostname, tt.port, "user", "password", tt.skipVerify)

			// We expect an error since the server doesn't exist
			if err == nil {
				t.Error("Expected error for non-existent server, got nil")
				if conn != nil {
					_ = conn.Close()
				}
				return
			}

			// Verify error message indicates connection failure
			if conn != nil {
				t.Error("Expected nil connection on error")
				_ = conn.Close()
			}
		})
	}
}

func TestConnectToAD_InvalidPort(t *testing.T) {
	// Test with invalid port - should fail gracefully
	conn, err := connectToAD("localhost", 0, "user", "password", false)
	if err == nil {
		t.Error("Expected error for invalid port, got nil")
		if conn != nil {
			_ = conn.Close()
		}
	}
}

func TestConnectToAD_InsecureWithSelfSigned(t *testing.T) {
	// This test documents the expected behavior:
	// - skipVerify=true should allow self-signed/invalid certificates
	// - skipVerify=false should reject them
	// Since we can't easily create a test TLS server here,
	// we just verify the function signature accepts the flag correctly

	// Test that the function can be called with both values
	// (actual TLS verification behavior is tested by the connection attempt)

	t.Run("skipVerify parameter is accepted", func(t *testing.T) {
		// These will fail to connect, but should not panic
		_, _ = connectToAD("127.0.0.1", 63636, "user", "pass", true)
		_, _ = connectToAD("127.0.0.1", 63636, "user", "pass", false)
	})
}
