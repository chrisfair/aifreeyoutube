#!/usr/bin/env bash
# ==============================================================================
# Build Windows Binary & Packaging Bundle (.exe / .zip)
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PACKAGE_NAME="aifreeyoutube"
VERSION="1.0.0"
OUTPUT_DIR="${ROOT_DIR}/build/dist"
STAGING_DIR="${ROOT_DIR}/build/windows-staging"

echo "==> Building Windows (amd64) distribution for ${PACKAGE_NAME} v${VERSION}..."
mkdir -p "${OUTPUT_DIR}"
rm -rf "${STAGING_DIR}"
mkdir -p "${STAGING_DIR}"

cd "${ROOT_DIR}"

# 1. Build Vue 3 production assets
echo "==> Building frontend assets..."
cd frontend && npm run build
cd "${ROOT_DIR}"

# 2. Compile standalone Windows web server (Zero CGO dependency)
echo "==> Compiling Windows server binary (aifreeyoutube-server.exe)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o "${STAGING_DIR}/aifreeyoutube-server.exe" ./cmd/server

# 3. Attempt Wails Windows desktop compilation if mingw cross-compiler is present
WAILS=$(which wails 2>/dev/null || echo "${HOME}/go/bin/wails")
if which x86_64-w64-mingw32-gcc >/dev/null 2>&1; then
    echo "==> Cross-compiling Wails native Windows desktop executable..."
    ${WAILS} build -platform windows/amd64 -clean
    if [[ -f "build/bin/aifreeyoutube.exe" ]]; then
        cp "build/bin/aifreeyoutube.exe" "${STAGING_DIR}/"
    fi
else
    echo "==> Note: 'x86_64-w64-mingw32-gcc' not found on host."
    echo "    Providing native standalone Windows server executable (aifreeyoutube-server.exe)."
    echo "    (To cross-compile the Wails WebView2 desktop binary on Arch: sudo pacman -S mingw-w64-gcc)"
fi

# 4. Copy icons and readme to staging
cp "${ROOT_DIR}/build/appicon.png" "${STAGING_DIR}/"
cp "${ROOT_DIR}/README.md" "${STAGING_DIR}/"

# 5. Create launch batch script for Windows users
cat << 'EOF' > "${STAGING_DIR}/Start-AI-Free-YouTube.bat"
@echo off
title AI-Free YouTube Web App
echo ==================================================
echo       AI-Free YouTube - Windows Local Server
echo ==================================================
echo Starting server on http://localhost:8484 ...
start "" http://localhost:8484
aifreeyoutube-server.exe
pause
EOF

# 6. Package into ZIP distribution
ZIP_FILE="${OUTPUT_DIR}/${PACKAGE_NAME}-v${VERSION}-windows-amd64.zip"
echo "==> Creating Windows ZIP release bundle..."
(
    cd "${STAGING_DIR}"
    if which zip >/dev/null 2>&1; then
        zip -r -q "${ZIP_FILE}" .
    else
        python3 -m zipfile -c "${ZIP_FILE}" *
    fi
)

echo "==> Windows distribution bundle created successfully:"
echo "    ${ZIP_FILE}"
ls -lh "${ZIP_FILE}"
