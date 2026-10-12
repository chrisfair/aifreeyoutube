#!/usr/bin/env bash
# ==============================================================================
# Build Snap Package for AI-Free YouTube
# ==============================================================================
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SNAP_DIR="${ROOT_DIR}/packaging/snap"
OUTPUT_DIR="${ROOT_DIR}/build/dist"

echo "==> Preparing Snap package build..."

# 1. Ensure binaries and assets exist
cd "${ROOT_DIR}"
if [[ ! -f "build/bin/aifreeyoutube" ]]; then
    make desktop
fi
if [[ ! -f "bin/aifreeyoutube-server" ]]; then
    make server
fi

# 2. Check for snapcraft
if ! which snapcraft >/dev/null 2>&1; then
    echo "==> WARNING: 'snapcraft' is not installed on this system."
    echo "    On Arch Linux / CachyOS: snap install snapcraft --classic"
    echo "    On Ubuntu / Debian:       sudo snap install snapcraft --classic"
    echo ""
    echo "    The snapcraft recipe is configured at: ${SNAP_DIR}/snapcraft.yaml"
    echo "    To build when snapcraft is available, run:"
    echo "      cd ${SNAP_DIR} && snapcraft"
    exit 0
fi

# 3. Execute snapcraft
mkdir -p "${OUTPUT_DIR}"
cd "${SNAP_DIR}"
snapcraft
find "${SNAP_DIR}" -name "*.snap" -exec cp {} "${OUTPUT_DIR}/" \;

echo "==> Snap package generated successfully in ${OUTPUT_DIR}:"
ls -lh "${OUTPUT_DIR}"/*.snap
