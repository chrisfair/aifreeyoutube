#!/usr/bin/env bash
# ==============================================================================
# Unified Package Builder for AI-Free YouTube
# ==============================================================================
# Builds packages for Arch, Debian, RedHat/RPM, Windows, macOS, AppImage, Flatpak, Snap
#
# Usage:
#   ./packaging/scripts/build-all.sh [target]
#
# Targets:
#   all (default)  - Build all possible packages on current system
#   deb            - Debian/Ubuntu (.deb)
#   rpm            - Fedora/RHEL/openSUSE (.rpm)
#   windows        - Windows 64-bit release bundle (.zip / .exe)
#   macos          - macOS Apple Silicon & Intel bundles (.tar.gz / Homebrew)
#   appimage       - Linux universal AppImage
#   flatpak        - Flatpak bundle
#   snap           - Canonical Snap package
#   arch           - Arch Linux package (.pkg.tar.zst)
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TARGET="${1:-all}"
OUTPUT_DIR="${ROOT_DIR}/build/dist"
mkdir -p "${OUTPUT_DIR}"

cd "${ROOT_DIR}"

log_step() {
    echo ""
    echo "=============================================================================="
    echo "  $1"
    echo "=============================================================================="
}

build_arch() {
    log_step "Building Arch Linux Package (.pkg.tar.xz / .pkg.tar.zst)"
    if which makepkg >/dev/null 2>&1; then
        BUILDDIR=/tmp/makepkg-aifreeyoutube makepkg -f -d
        find . -maxdepth 1 -name "*.pkg.tar.*" -exec mv -f {} "${OUTPUT_DIR}/" \;
        echo "==> Arch Linux package moved to ${OUTPUT_DIR}"
    else
        echo "==> makepkg not available on this system."
    fi
}

build_deb() {
    log_step "Building Debian / Ubuntu Package (.deb)"
    "${ROOT_DIR}/packaging/debian/build-deb.sh"
}

build_rpm() {
    log_step "Building Red Hat / Fedora / openSUSE Package (.rpm)"
    "${ROOT_DIR}/packaging/rpm/build-rpm.sh" || true
}

build_windows() {
    log_step "Building Windows Distribution (.exe / .zip / Chocolatey)"
    "${ROOT_DIR}/packaging/windows/build-windows.sh"
}

build_macos() {
    log_step "Building macOS Distribution & Homebrew Tarballs"
    "${ROOT_DIR}/packaging/macos/build-macos.sh"
}

build_appimage() {
    log_step "Building Linux Universal AppImage"
    "${ROOT_DIR}/packaging/appimage/build-appimage.sh" || true
}

build_flatpak() {
    log_step "Preparing Flatpak Package"
    "${ROOT_DIR}/packaging/flatpak/build-flatpak.sh" || true
}

build_snap() {
    log_step "Preparing Snap Package"
    "${ROOT_DIR}/packaging/snap/build-snap.sh" || true
}

case "${TARGET}" in
    arch)     build_arch ;;
    deb)      build_deb ;;
    rpm)      build_rpm ;;
    windows)  build_windows ;;
    macos)    build_macos ;;
    appimage) build_appimage ;;
    flatpak)  build_flatpak ;;
    snap)     build_snap ;;
    all)
        build_arch
        build_deb
        build_rpm
        build_windows
        build_macos
        build_appimage
        build_flatpak
        build_snap
        ;;
    *)
        echo "Unknown target: ${TARGET}"
        echo "Available targets: all, arch, deb, rpm, windows, macos, appimage, flatpak, snap"
        exit 1
        ;;
esac

echo ""
echo "=============================================================================="
echo "  Packaging Complete! Summary of Generated Artifacts in build/dist/:"
echo "=============================================================================="
ls -lh "${OUTPUT_DIR}"/ || true
