#!/usr/bin/env bash
# ==============================================================================
# AI-Free YouTube: AUR Publishing Script
# ==============================================================================
# Pushes PKGBUILD and .SRCINFO to aur.archlinux.org using ~/.ssh/aur
# All packaging manifests stay version-controlled within this repository.
# ==============================================================================

set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
AUR_URL="ssh://aur@aur.archlinux.org/aifreeyoutube.git"

echo "==> Regenerating .SRCINFO from PKGBUILD..."
cd "${REPO_DIR}/packaging/aur"
makepkg --printsrcinfo > .SRCINFO
cp PKGBUILD .SRCINFO "${REPO_DIR}/"

TMP_AUR="$(mktemp -d /tmp/aur-publish.XXXXXX)"
trap 'rm -rf "${TMP_AUR}"' EXIT

echo "==> Cloning or initializing AUR package base..."
if ! git clone "${AUR_URL}" "${TMP_AUR}" 2>/dev/null; then
    mkdir -p "${TMP_AUR}"
    cd "${TMP_AUR}"
    git init -b master
    git remote add origin "${AUR_URL}"
fi

cd "${TMP_AUR}"
git checkout -B master 2>/dev/null || true
cp "${REPO_DIR}/packaging/aur/PKGBUILD" "${TMP_AUR}/PKGBUILD"
cp "${REPO_DIR}/packaging/aur/.SRCINFO" "${TMP_AUR}/.SRCINFO"

git add PKGBUILD .SRCINFO
if git diff --staged --quiet; then
    echo "✔ AUR package is already up to date. Nothing to push."
else
    git commit -m "chore(aur): release aifreeyoutube 1.0.0-1"
    echo "==> Pushing to aur.archlinux.org (using ~/.ssh/aur)..."
    git push origin master
    echo "✔ Successfully published to AUR!"
fi
