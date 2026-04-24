# 🧪 Test Suite

The project includes a comprehensive test suite located in `src/adtest_test.go`.

---

## Running Tests

```bash
cd src
go test -v ./...
```

Or use the build script which runs linting and tests before building:

```bash
./build.sh
```

The build script performs:
1. **Linting** with golangci-lint (requires v2 config format)
2. **Testing** with race detection enabled (`CGO_ENABLED=1`)
3. **Building** static binary with optimization flags (`CGO_ENABLED=0`)

### Running Specific Tests

Run tests matching a pattern:

```bash
go test -v -run "TestLoadConfig" ./...
```

### Running with Race Detection

```bash
go test -v -race ./...
```

### Running with Coverage

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## Test Coverage

The test suite covers the following areas:

| Test Name | Description |
| --------- | ----------- |
| `TestLoadConfig` | Verifies loading a valid JSON configuration file |
| `TestLoadConfigNonExistent` | Ensures proper error handling for missing files |
| `TestLoadConfigInvalidJSON` | Validates error handling for malformed JSON |
| `TestGenerateConfigFile` | Tests config file generation and round-trip loading |
| `TestConfigValidation` | Table-driven tests for required field validation |
| `TestConfig_Validate_PortRange` | Validates port number range (1-65535) |
| `TestConnectToAD_InsecureFlag` | Tests both secure and insecure TLS modes |
| `TestConnectToAD_InvalidPort` | Verifies graceful handling of invalid ports |
| `TestConnectToAD_InsecureWithSelfSigned` | Documents expected behavior for certificate verification |

---

## Test Categories

### Configuration Tests

Tests for loading, parsing, and validating JSON configuration files:

- **Valid config loading**: Ensures all fields are correctly parsed
- **Missing file handling**: Returns appropriate error for non-existent files
- **Invalid JSON handling**: Gracefully handles malformed JSON input
- **Config generation**: Verifies example config can be created and re-loaded
- **Config validation**: Uses `Config.Validate()` method for proper encapsulation
- **Port range validation**: Ensures port is within valid range (1-65535)

### Connection Tests

Tests for the LDAP connection functionality:

- **Insecure flag**: Validates both `skipVerify=true` and `skipVerify=false` modes
- **Invalid port**: Ensures graceful failure with invalid connection parameters
- **Self-signed certificates**: Documents behavior for certificate verification

---

## Test Design Principles

- **Table-driven tests**: Used for configuration validation to cover multiple scenarios efficiently
- **Temporary files**: Tests create and clean up temp files to avoid side effects
- **Error path coverage**: Both success and failure cases are tested
- **No external dependencies**: LDAP connection tests verify behavior without requiring a live server
- **Isolated tests**: Each test is independent and can run in any order

---

## Adding New Tests

When adding new functionality, follow these guidelines:

1. **Name tests descriptively**: Use `TestFunctionName_Scenario` pattern
2. **Use table-driven tests**: For multiple input scenarios
3. **Clean up resources**: Use `defer` for temp files and connections
4. **Test error paths**: Don't just test the happy path
5. **Document edge cases**: Use comments to explain non-obvious test cases

### Example: Table-Driven Test

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "result", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := NewFeature(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

---

## Linting Configuration

The project uses **golangci-lint v2** for static analysis. Configuration is in `.golangci.yml`:

```yaml
version: "2"  # Required for golangci-lint v2.x

linters:
  default: standard
  enable:
    - errcheck      # Check for unchecked errors
    - govet         # Go vet checks
    - staticcheck   # Static analysis
    - unused        # Find unused code
    - misspell      # Spelling mistakes
    - gocritic      # Code quality hints

  settings:        # Note: in v2, settings are nested under linters
    gocritic:
      enabled-tags:
        - diagnostic
        - style
        - performance
```

### Installing golangci-lint

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### Common Linting Issues

| Issue | Fix |
|-------|-----|
| `errcheck: Error return value not checked` | Use `_ = f.Close()` for intentional ignores |
| `exitAfterDefer: log.Fatal will exit, defer won't run` | Use `log.Println()` + `return` instead |
| `staticcheck: deprecated API` | Use recommended replacement (e.g., `ldap.DialURL`) |
| `gocritic: paramTypeCombine` | Use `func(a, b string)` instead of `func(a string, b string)` |

---

## Continuous Integration

The `build.sh` script integrates tests into the build pipeline:

1. Tests run **before** compilation
2. Build fails immediately if any test fails
3. Prevents broken binaries from being produced

```bash
./build.sh
```

Output:
```
Running linter...
0 issues.
✅ Linting passed
Running tests with race detection...
=== RUN   TestLoadConfig
--- PASS: TestLoadConfig (0.00s)
...
✅ Tests passed
Building adtest binary...
✅ Build complete: bin/adtest
```
