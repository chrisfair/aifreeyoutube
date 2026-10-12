#!/usr/bin/env bash
# ==============================================================================
# Build Debian / Ubuntu (.deb) Package for AI-Free YouTube
# ==============================================================================
# Creates a standard Debian package installable via:
#   sudo dpkg -i aifreeyoutube_1.0.0_amd64.deb
#   sudo apt-get install -f
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PACKAGE_NAME="aifreeyoutube"
VERSION="1.0.0"
ARCH="amd64"
DEB_DIR="${ROOT_DIR}/build/deb-staging"
OUTPUT_DIR="${ROOT_DIR}/build/dist"

echo "==> Building Debian package for ${PACKAGE_NAME} v${VERSION} (${ARCH})..."

# 1. Ensure binaries and assets are built
cd "${ROOT_DIR}"
if [[ ! -f "build/bin/aifreeyoutube" ]]; then
    echo "==> Compiling desktop application..."
    make desktop
fi

if [[ ! -f "bin/aifreeyoutube-server" ]]; then
    echo "==> Compiling standalone web server..."
    make server
fi

# 2. Prepare staging tree
rm -rf "${DEB_DIR}"
mkdir -p "${DEB_DIR}/DEBIAN"
mkdir -p "${DEB_DIR}/usr/bin"
mkdir -p "${DEB_DIR}/usr/share/applications"
mkdir -p "${DEB_DIR}/usr/share/icons/hicolor/scalable/apps"
mkdir -p "${DEB_DIR}/usr/share/icons/hicolor/512x512/apps"
mkdir -p "${DEB_DIR}/lib/systemd/system"
mkdir -p "${OUTPUT_DIR}"

# 3. Copy control files
cp "${ROOT_DIR}/packaging/debian/control" "${DEB_DIR}/DEBIAN/control"
cp "${ROOT_DIR}/packaging/debian/postinst" "${DEB_DIR}/DEBIAN/postinst"
cp "${ROOT_DIR}/packaging/debian/postrm" "${DEB_DIR}/DEBIAN/postrm"
chmod 755 "${DEB_DIR}/DEBIAN/postinst" "${DEB_DIR}/DEBIAN/postrm"

# 4. Copy application binaries
install -m 755 "${ROOT_DIR}/build/bin/aifreeyoutube" "${DEB_DIR}/usr/bin/aifreeyoutube"
install -m 755 "${ROOT_DIR}/bin/aifreeyoutube-server" "${DEB_DIR}/usr/bin/aifreeyoutube-server"

# 5. Copy desktop integration assets
install -m 644 "${ROOT_DIR}/build/linux/aifreeyoutube.desktop" "${DEB_DIR}/usr/share/applications/aifreeyoutube.desktop"
install -m 644 "${ROOT_DIR}/build/appicon.svg" "${DEB_DIR}/usr/share/icons/hicolor/scalable/apps/aifreeyoutube.svg"
install -m 644 "${ROOT_DIR}/build/appicon.png" "${DEB_DIR}/usr/share/icons/hicolor/512x512/apps/aifreeyoutube.png"

# 6. Copy systemd service
install -m 644 "${ROOT_DIR}/packaging/systemd/aifreeyoutube.service" "${DEB_DIR}/lib/systemd/system/aifreeyoutube.service"

# 7. Build package using dpkg-deb or ar fallback
OUTPUT_FILE="${OUTPUT_DIR}/${PACKAGE_NAME}_${VERSION}_${ARCH}.deb"
if which dpkg-deb >/dev/null 2>&1; then
    dpkg-deb --build --root-owner-group "${DEB_DIR}" "${OUTPUT_FILE}"
else
    echo "==> dpkg-deb not found; using ar + tar fallback..."
    (
        cd "${DEB_DIR}"
        echo "2.0" > debian-binary
        tar -czf control.tar.gz -C DEBIAN .
        tar -czf data.tar.gz --exclude=DEBIAN --exclude=debian-binary --exclude=control.tar.gz .
        ar rcs "${OUTPUT_FILE}" debian-binary control.tar.gz data.tar.gz
    )
fi

echo "==> Debian package generated successfully:"
echo "    ${OUTPUT_FILE}"
ls -lh "${OUTPUT_FILE}"
