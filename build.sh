#!/usr/bin/env bash
set -euo pipefail

# Ensure Go bin directory is in PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Force rebuild flag
FORCE=""
if [[ "${1:-}" == "1" ]]; then
    FORCE="-a"
fi

# Run linting (if available and not skipped)
cd src
if [[ "${SKIP_LINT:-}" == "1" ]]; then
    echo "⚠️  Skipping linting (SKIP_LINT=1)"
elif command -v golangci-lint &> /dev/null; then
    echo "Running linter..."
    if ! golangci-lint run ./...; then
        echo "❌ Linting failed"
        exit 1
    fi
    echo "✅ Linting passed"
else
    echo "⚠️  golangci-lint not installed, skipping linting"
fi

# Run tests with race detection (requires CGO)
echo "Running tests with race detection..."
export CGO_ENABLED=1

if ! go test -v -race ./...; then
    echo "❌ Tests failed"
    exit 1
fi

echo "✅ Tests passed"

# Build configuration (static binary, no CGO)
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Build with optimization flags
echo "Building adtest binary..."

if ! go build -trimpath -ldflags="-s -w" $FORCE -o ../bin/adtest .; then
    echo "❌ Build failed"
    exit 1
fi

cd ..

echo "✅ Build complete: bin/adtest"
ls -lh bin/adtest
