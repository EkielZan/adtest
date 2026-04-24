#!/usr/bin/env bash
set -euo pipefail

# Force rebuild flag
FORCE=""
if [[ "${1:-}" == "1" ]]; then
    FORCE="-a"
fi

# Build configuration
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Run linting (if available)
cd src
echo "Running linter..."
if command -v golangci-lint &> /dev/null; then
    if ! golangci-lint run ./...; then
        echo "❌ Linting failed"
        exit 1
    fi
    echo "✅ Linting passed"
else
    echo "⚠️  golangci-lint not installed, skipping linting"
fi

# Run tests with race detection
echo "Running tests with race detection..."

if ! go test -v -race ./...; then
    echo "❌ Tests failed"
    exit 1
fi

echo "✅ Tests passed"

# Build with optimization flags
echo "Building adtest binary..."

if ! go build -trimpath -ldflags="-s -w" $FORCE -o ../bin/adtest .; then
    echo "❌ Build failed"
    exit 1
fi

cd ..

echo "✅ Build complete: bin/adtest"
ls -lh bin/adtest
