package main

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
