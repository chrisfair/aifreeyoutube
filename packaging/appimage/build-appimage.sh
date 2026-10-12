#!/usr/bin/env bash
# ==============================================================================
# Build Linux Universal AppImage for AI-Free YouTube
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PACKAGE_NAME="aifreeyoutube"
VERSION="1.0.0"
APPDIR="${ROOT_DIR}/build/AppDir"
OUTPUT_DIR="${ROOT_DIR}/build/dist"

echo "==> Preparing AppImage structure for ${PACKAGE_NAME} v${VERSION}..."

# 1. Build desktop binary
cd "${ROOT_DIR}"
if [[ ! -f "build/bin/aifreeyoutube" ]]; then
    make desktop
fi

# 2. Setup AppDir directory tree
rm -rf "${APPDIR}"
mkdir -p "${APPDIR}/usr/bin"
mkdir -p "${APPDIR}/usr/share/applications"
mkdir -p "${APPDIR}/usr/share/icons/hicolor/scalable/apps"
mkdir -p "${APPDIR}/usr/share/icons/hicolor/512x512/apps"
mkdir -p "${OUTPUT_DIR}"

# 3. Copy application binary
install -m 755 "${ROOT_DIR}/build/bin/aifreeyoutube" "${APPDIR}/usr/bin/aifreeyoutube"

# 4. Copy icons and desktop entry
install -m 644 "${ROOT_DIR}/build/linux/aifreeyoutube.desktop" "${APPDIR}/usr/share/applications/aifreeyoutube.desktop"
install -m 644 "${ROOT_DIR}/build/linux/aifreeyoutube.desktop" "${APPDIR}/aifreeyoutube.desktop"
install -m 644 "${ROOT_DIR}/build/appicon.svg" "${APPDIR}/usr/share/icons/hicolor/scalable/apps/aifreeyoutube.svg"
install -m 644 "${ROOT_DIR}/build/appicon.png" "${APPDIR}/usr/share/icons/hicolor/512x512/apps/aifreeyoutube.png"
install -m 644 "${ROOT_DIR}/build/appicon.png" "${APPDIR}/aifreeyoutube.png"
install -m 644 "${ROOT_DIR}/build/appicon.png" "${APPDIR}/.DirIcon"

# 5. Copy AppRun
install -m 755 "${ROOT_DIR}/packaging/appimage/AppRun" "${APPDIR}/AppRun"

# 6. Build AppImage using appimagetool if available, or fetch standalone appimagetool
APPIMAGE_BIN="${OUTPUT_DIR}/${PACKAGE_NAME}-${VERSION}-x86_64.AppImage"

if which appimagetool >/dev/null 2>&1; then
    echo "==> Building AppImage using system appimagetool..."
    ARCH=x86_64 appimagetool "${APPDIR}" "${APPIMAGE_BIN}"
else
    echo "==> 'appimagetool' not found in PATH."
    TOOL_PATH="${ROOT_DIR}/build/appimagetool-x86_64.AppImage"
    if [[ ! -f "${TOOL_PATH}" ]]; then
        echo "==> Fetching standalone appimagetool from GitHub releases..."
        curl -sSL -o "${TOOL_PATH}" "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage" || true
        chmod +x "${TOOL_PATH}" || true
    fi

    if [[ -x "${TOOL_PATH}" ]]; then
        echo "==> Packaging with local appimagetool..."
        ARCH=x86_64 "${TOOL_PATH}" --appimage-extract-and-run "${APPDIR}" "${APPIMAGE_BIN}" || true
    fi
fi

if [[ -f "${APPIMAGE_BIN}" ]]; then
    echo "==> AppImage created successfully: ${APPIMAGE_BIN}"
    ls -lh "${APPIMAGE_BIN}"
else
    echo "==> AppDir assembled at: ${APPDIR}"
    echo "    To convert to .AppImage manually, run: appimagetool build/AppDir"
fi
