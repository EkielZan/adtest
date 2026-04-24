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

# Build with optimization flags
cd src
echo "Building adtest binary..."

if ! go build -trimpath -ldflags="-s -w" $FORCE -o ../bin/adtest .; then
    echo "❌ Build failed"
    exit 1
fi

cd ..

echo "✅ Build complete: bin/adtest"
ls -lh bin/adtest
