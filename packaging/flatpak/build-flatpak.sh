#!/usr/bin/env bash
# ==============================================================================
# Build Flatpak Bundle for AI-Free YouTube
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
MANIFEST="${ROOT_DIR}/packaging/flatpak/com.aifreeyoutube.app.yml"
BUILD_DIR="${ROOT_DIR}/build/flatpak-build"
REPO_DIR="${ROOT_DIR}/build/flatpak-repo"
OUTPUT_DIR="${ROOT_DIR}/build/dist"

echo "==> Preparing Flatpak build for com.aifreeyoutube.app..."

# 1. Ensure binaries and assets exist
cd "${ROOT_DIR}"
if [[ ! -f "build/bin/aifreeyoutube" ]]; then
    make desktop
fi
if [[ ! -f "bin/aifreeyoutube-server" ]]; then
    make server
fi

# 2. Check for flatpak-builder
if ! which flatpak-builder >/dev/null 2>&1; then
    echo "==> WARNING: 'flatpak-builder' is not installed on this system."
    echo "    On Arch Linux: sudo pacman -S flatpak-builder"
    echo "    On Fedora:     sudo dnf install flatpak-builder"
    echo "    On Ubuntu:     sudo apt install flatpak-builder"
    echo ""
    echo "    The manifest has been generated at: ${MANIFEST}"
    echo "    When flatpak-builder is available, run:"
    echo "      flatpak-builder --user --install --force-clean ${BUILD_DIR} ${MANIFEST}"
    exit 0
fi

# 3. Build Flatpak
mkdir -p "${OUTPUT_DIR}"
flatpak-builder --force-clean --repo="${REPO_DIR}" "${BUILD_DIR}" "${MANIFEST}"
flatpak build-bundle "${REPO_DIR}" "${OUTPUT_DIR}/aifreeyoutube.flatpak" com.aifreeyoutube.app

echo "==> Flatpak bundle generated successfully:"
echo "    ${OUTPUT_DIR}/aifreeyoutube.flatpak"
ls -lh "${OUTPUT_DIR}/aifreeyoutube.flatpak"
