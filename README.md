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
## ⚙️ Configuration file

`adtest` can be configured using a **JSON configuration file**, supplied via the `-config`
command-line flag. This allows connection details to be stored centrally and overridden
when necessary using CLI flags.

### Example configuration file (`config.json`)

```json
{
  "hostname": "SA1000000047.ad.ing.net",
  "port": 636,
  "baseDN": "DC=ad,DC=ing,DC=net",
  "bindUser": "AD\\service_account"
}
````

### Configuration fields explained

| Field      | Description                                                            |
| ---------- | ---------------------------------------------------------------------- |
| `hostname` | Fully qualified domain name (FQDN) of the Active Directory server      |
| `port`     | LDAP port number (typically `636` for LDAPS)                           |
| `baseDN`   | Base Distinguished Name used as the root for LDAP searches             |
| `bindUser` | User account used to authenticate to Active Directory (`DOMAIN\\user`) |

⚠️ All configuration fields are required unless overridden via command-line flags.

***

### Using a configuration file

```bash
./adtest -config config.json -sam testuser
```

When the `-config` flag is provided, the configuration file is loaded before processing
command-line overrides.

***

## 🔀 Configuration precedence

Configuration values are resolved using the following order (highest priority first):

1.  **Command-line flags**
2.  **Configuration file**
3.  **Built-in defaults**

This allows reuse of the same configuration file across environments while still enabling
runtime overrides.

### Example: override bind user only

```bash
./adtest \
  -config config.json \
  -b "AD\\override_user" \
  -sam testuser
```

In this case:

*   `hostname`, `port`, and `baseDN` come from the config file
*   `bindUser` is overridden via the command line

***

## 🎛️ Command-line flags

Below is the list of all supported command-line flags.

### Required flags

| Flag   | Description                                |
| ------ | ------------------------------------------ |
| `-sam` | `sAMAccountName` of the user to search for |

***

### Optional flags

| Flag      | Description                                                         |
| --------- | ------------------------------------------------------------------- |
| `-config` | Path to a JSON configuration file                                   |
| `-b`      | Bind user used to authenticate to Active Directory (`DOMAIN\\user`) |
| `-json`   | Output the result in JSON format instead of human-readable text     |

***

### Example usages

#### Use configuration file only

```bash
./adtest -config config.json -sam testuser
```

#### Override bind user from CLI

```bash
./adtest -config config.json -b "AD\\svc_alt" -sam testuser
```

#### Enable JSON output for scripting

```bash
./adtest -config config.json -sam testuser -json
```

***

## 🔐 Password handling

Passwords are **never** stored in configuration files and are **never** accepted via
command-line flags.

The password is obtained using one of the following methods:

1.  From the `LDAP_PASSWORD` environment variable (non-interactive usage)
2.  Secure terminal prompt (interactive usage, input not echoed)

This design prevents accidental disclosure via logs, shell history, or process listings.

***

## 🤖 Project origin

This project was created with the assistance of an AI-based coding helper, which was used
to accelerate development, improve code quality, and refine documentation.  
All design decisions, implementation choices, and final validation remain the
responsibility of the project author.

```
