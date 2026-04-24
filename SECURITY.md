# Security Policy

## Security Improvements (2026-04-23)

This project has undergone a comprehensive security review and the following improvements have been implemented:

### ✅ Fixed Vulnerabilities

1. **TLS Certificate Validation** (CRITICAL)
   - **Issue**: Certificate validation was disabled (`InsecureSkipVerify: true`)
   - **Fix**: Enabled proper TLS certificate validation with ServerName verification
   - **Impact**: Prevents man-in-the-middle attacks on LDAPS connections

2. **LDAP Injection** (HIGH)
   - **Issue**: User input was not escaped in LDAP queries
   - **Fix**: All user input is now escaped using `ldap.EscapeFilter()`
   - **Impact**: Prevents LDAP injection attacks

3. **Updated Dependencies** (HIGH)
   - **Issue**: Go 1.13 and dependencies 3-4 years outdated with known CVEs
   - **Fix**: Updated to Go 1.22+ and all dependencies to latest versions
   - **Impact**: Patches multiple security vulnerabilities

### Security Features

- **TLS 1.2+ Only**: Minimum TLS version enforced
- **Input Validation**: All user inputs are escaped before LDAP queries
- **Secure Dependencies**: All dependencies updated to latest secure versions
- **No Debug Symbols**: Production binaries built without debug information

### Best Practices

When using this tool:

1. **Always use LDAPS** (port 636) - never plain LDAP
2. **Use environment variables** for passwords (`LDAP_PASSWORD`)
3. **Validate certificates** - do not disable certificate validation
4. **Use service accounts** with minimal required permissions
5. **Rotate credentials** regularly

## Reporting a Vulnerability

If you discover a security vulnerability, please email the maintainers directly rather than opening a public issue.

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.

## Security Updates

- 2026-04-23: Major security fixes applied (TLS validation, LDAP injection, dependency updates)
