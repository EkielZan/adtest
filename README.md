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
- Script-friendly exit codes- Optional insecure mode for self-signed/untrusted certificates
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
  "hostname": "hostname.company.com",
  "port": 636,
  "baseDN": "DC=ad,DC=company,DC=com",
  "bindUser": "AD\\service_account"
}
```

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

| Flag        | Description                                                         |
| ----------- | ------------------------------------------------------------------- |
| `-config`   | Path to a JSON configuration file                                   |
| `-b`        | Bind user used to authenticate to Active Directory (`DOMAIN\\user`) |
| `-json`     | Output the result in JSON format instead of human-readable text     |
| `-json-gen` | Generate an example JSON configuration file (`example.json`)        |
| `-insecure` | Skip TLS certificate verification (use with caution)                |

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

#### Generate example configuration file

```bash
./adtest -json-gen
```

This creates an `example.json` file with default configuration values.

#### Connect with insecure mode (self-signed certificates)

```bash
./adtest -config config.json -sam testuser -insecure
```

***

## 🔓 TLS Certificate Verification

By default, `adtest` enforces strict TLS certificate verification with a minimum of TLS 1.2.
This ensures secure communication with your Active Directory server.

However, in development or test environments where self-signed or untrusted certificates
are used, you can disable certificate verification using the `-insecure` flag:

```bash
./adtest -config config.json -sam testuser -insecure
```

⚠️ **Security Warning**: Using `-insecure` disables certificate validation, making
the connection vulnerable to man-in-the-middle attacks. **Never use this flag in
production environments.**

When `-insecure` is used, a warning message is displayed:
```
⚠️  WARNING: TLS certificate verification is disabled
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

## 🧪 Testing

The project includes a comprehensive test suite. Tests are integrated into the build
pipeline and run automatically before compilation.

**Quick start:**

```bash
cd src
go test -v ./...
```

Or use the build script:

```bash
./build.sh
```

📖 **For detailed testing documentation**, including test coverage, design principles,
and guidelines for adding new tests, see **[TESTING.md](TESTING.md)**.

***

## 🤖 Project origin

This project was created with the assistance of an AI-based coding helper, which was used
to accelerate development, improve code quality, and refine documentation.  
All design decisions, implementation choices, and final validation remain the
responsibility of the project author.

---

## 📝 AI Contribution Notes

The following components were generated or enhanced with AI assistance:

- **Documentation**: This README and inline code comments
- **Test cases**: Unit tests in `adtest_test.go`
- **Build pipeline**: `build.sh` with test integration and workflow improvements
- **Agent configuration**: `.github/agents/go-build-reviewer.agent.md`

**AI Version**: Claude Opus 4.5 (GitHub Copilot)
