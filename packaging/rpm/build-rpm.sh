#!/usr/bin/env bash
# ==============================================================================
# Build RPM Package for AI-Free YouTube (Fedora, RHEL, openSUSE)
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PACKAGE_NAME="aifreeyoutube"
VERSION="1.0.0"
RPM_TOPDIR="${ROOT_DIR}/build/rpm"
OUTPUT_DIR="${ROOT_DIR}/build/dist"

echo "==> Building RPM package for ${PACKAGE_NAME} v${VERSION}..."

# Check for rpmbuild tool
if ! which rpmbuild >/dev/null 2>&1; then
    echo "==> WARNING: 'rpmbuild' is not installed on this system."
    echo "    On Arch Linux, install it via: sudo pacman -S rpm-tools"
    echo "    On Fedora / RHEL: sudo dnf install rpm-build"
    echo "    On openSUSE: sudo zypper install rpm-build"
    exit 1
fi

# Ensure binaries and assets are built
cd "${ROOT_DIR}"
if [[ ! -f "build/bin/aifreeyoutube" ]]; then
    make desktop
fi
if [[ ! -f "bin/aifreeyoutube-server" ]]; then
    make server
fi

# Setup RPM directory layout
rm -rf "${RPM_TOPDIR}"
mkdir -p "${RPM_TOPDIR}/"{BUILD,RPMS,SOURCES,SPECS,SRPMS}
mkdir -p "${OUTPUT_DIR}"

# Stage sources into SOURCES directory for rpmbuild
cp "${ROOT_DIR}/build/bin/aifreeyoutube" "${RPM_TOPDIR}/SOURCES/"
cp "${ROOT_DIR}/bin/aifreeyoutube-server" "${RPM_TOPDIR}/SOURCES/"
cp "${ROOT_DIR}/build/linux/aifreeyoutube.desktop" "${RPM_TOPDIR}/SOURCES/"
cp "${ROOT_DIR}/build/appicon.svg" "${RPM_TOPDIR}/SOURCES/"
cp "${ROOT_DIR}/build/appicon.png" "${RPM_TOPDIR}/SOURCES/"
cp "${ROOT_DIR}/packaging/systemd/aifreeyoutube.service" "${RPM_TOPDIR}/SOURCES/"
cp "${ROOT_DIR}/packaging/rpm/aifreeyoutube.spec" "${RPM_TOPDIR}/SPECS/"

# Build the RPM
rpmbuild -bb \
    --define "_topdir ${RPM_TOPDIR}" \
    --define "_sourcedir ${RPM_TOPDIR}/SOURCES" \
    "${RPM_TOPDIR}/SPECS/aifreeyoutube.spec"

# Copy resulting RPM to dist directory
find "${RPM_TOPDIR}/RPMS" -name "*.rpm" -exec cp {} "${OUTPUT_DIR}/" \;
echo "==> RPM package built successfully in ${OUTPUT_DIR}:"
ls -lh "${OUTPUT_DIR}"/*.rpm
