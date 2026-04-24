package main

import "fmt"

type Config struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	BaseDN   string `json:"baseDN"`
	BindUser string `json:"bindUser"`
}

// Validate checks that all required configuration fields are present and valid.
func (c *Config) Validate() error {
	if c.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}
	if c.BaseDN == "" {
		return fmt.Errorf("baseDN is required")
	}
	if c.BindUser == "" {
		return fmt.Errorf("bindUser is required")
	}
	return nil
}

type UserResult struct {
	DN             string `json:"dn"`
	CN             string `json:"cn"`
	SAMAccountName string `json:"sAMAccountName"`
	DisplayName    string `json:"displayName,omitempty"`
	Mail           string `json:"mail,omitempty"`
}
