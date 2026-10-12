# Maintainer: Christopher Patrick Fair <christopherpatrickfair@gmail.com>
# Contributor: AI-Free YouTube Project Maintainers

pkgname=aifreeyoutube
pkgver=1.0.0
pkgrel=1
pkgdesc="Authentic human video portal filtering out AI slop, clickbait mills, and bots on YouTube"
arch=('x86_64')
url="https://github.com/christopherpatrickfair/aifreeyoutube"
license=('GPL-3.0-or-later')
depends=(
    'webkit2gtk-4.1'
    'gtk3'
    'glib2'
    'hicolor-icon-theme'
)
makedepends=(
    'go'
    'nodejs'
    'npm'
)
provides=('aifreeyoutube')
conflicts=('aifreeyoutube')
options=('!strip')

build() {
    cd "$startdir"
    
    # Ensure wails CLI is discoverable in GOPATH if not in system PATH
    export PATH="$PATH:$HOME/go/bin"
    
    echo "==> Building AI-Free YouTube Wails desktop application..."
    wails build -tags webkit2_41 -clean
}

package() {
    cd "$startdir"
    
    # Install application binary
    install -Dm755 "build/bin/aifreeyoutube" "$pkgdir/usr/bin/aifreeyoutube"
    
    # Install Freedesktop .desktop entry
    install -Dm644 "build/linux/aifreeyoutube.desktop" "$pkgdir/usr/share/applications/aifreeyoutube.desktop"
    
    # Install scalable vector icon
    install -Dm644 "build/appicon.svg" "$pkgdir/usr/share/icons/hicolor/scalable/apps/aifreeyoutube.svg"
    
    # Install high-resolution 512x512 raster icon
    install -Dm644 "build/appicon.png" "$pkgdir/usr/share/icons/hicolor/512x512/apps/aifreeyoutube.png"

    # Install license
    install -Dm644 "LICENSE" "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
}
