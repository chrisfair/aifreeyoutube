#!/usr/bin/env bash
# ==============================================================================
# Build macOS Binaries and Homebrew Tarballs (Apple Silicon & Intel)
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PACKAGE_NAME="aifreeyoutube"
VERSION="1.0.0"
OUTPUT_DIR="${ROOT_DIR}/build/dist"
STAGING_DIR="${ROOT_DIR}/build/macos-staging"

echo "==> Building macOS release bundles for ${PACKAGE_NAME} v${VERSION}..."
mkdir -p "${OUTPUT_DIR}"

cd "${ROOT_DIR}"

# 1. Build Vue 3 production frontend
cd frontend && npm run build
cd "${ROOT_DIR}"

# 2. Build Apple Silicon (ARM64)
echo "==> Compiling macOS ARM64 binary..."
rm -rf "${STAGING_DIR}-arm64"
mkdir -p "${STAGING_DIR}-arm64"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o "${STAGING_DIR}-arm64/aifreeyoutube-server" ./cmd/server
cp README.md "${STAGING_DIR}-arm64/"
tar -czf "${OUTPUT_DIR}/${PACKAGE_NAME}-v${VERSION}-darwin-arm64.tar.gz" -C "${STAGING_DIR}-arm64" .

# 3. Build Intel Mac (AMD64)
echo "==> Compiling macOS AMD64 binary..."
rm -rf "${STAGING_DIR}-amd64"
mkdir -p "${STAGING_DIR}-amd64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o "${STAGING_DIR}-amd64/aifreeyoutube-server" ./cmd/server
cp README.md "${STAGING_DIR}-amd64/"
tar -czf "${OUTPUT_DIR}/${PACKAGE_NAME}-v${VERSION}-darwin-amd64.tar.gz" -C "${STAGING_DIR}-amd64" .

echo "==> macOS release bundles created successfully:"
ls -lh "${OUTPUT_DIR}"/*darwin*.tar.gz
