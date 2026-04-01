# adtest

`adtest` is a command-line utility written in Go that queries Microsoft Active
Directory over LDAPS (LDAP over TLS).  
It allows you to authenticate with a bind account and search for a user by
`sAMAccountName`.

The tool is designed to be:
- Script-friendly
- Secure by default
- Easy to configure via config files and CLI flags
- Comparable in usage to `ldapsearch`

---

## ✅ Features

- Secure LDAPS (TLS) connection to Active Directory
- Search users by `sAMAccountName`
- Configuration file support (JSON)
- Command-line flags override configuration file values
- Secure password handling (no hardcoded secrets, no echo)
- Optional JSON output for automation
- Clear handling of “user not found” cases
- Script-friendly exit codes

---

## 📦 Build

```bash
go build -o adtest

```

---

## 🤖 Project origin

This project was created with the assistance of an AI-based coding helper, which was used to accelerate development, improve code quality, and refine documentation.  
All design decisions, implementation choices, and final validation remain the responsibility of the project author.