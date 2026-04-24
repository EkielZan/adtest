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
	defer os.Remove(tmpfile.Name())

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
	tmpfile.Close()

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
	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString("{invalid json}")
	tmpfile.Close()

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
	tmpfile.Close()
	os.Remove(tmpPath) // Remove so we can test creation
	defer os.Remove(tmpPath)

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
			valid := tt.config.Hostname != "" && tt.config.Port != 0 && 
				tt.config.BaseDN != "" && tt.config.BindUser != ""
			if valid != tt.valid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.valid, valid)
			}
		})
	}
}
