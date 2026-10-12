# ==============================================================================
# AI-Free YouTube: Universal Cross-Platform Packaging & Desktop Engine
# ==============================================================================
# Supported targets: Arch Linux, Debian/Ubuntu, Red Hat/RPM, Windows/Chocolatey,
# macOS/Homebrew, Linux Universal AppImage, Flatpak, Canonical Snap.
# ==============================================================================

WAILS := $(shell which wails 2>/dev/null || echo $(HOME)/go/bin/wails)
PREFIX ?= $(HOME)/.local

.PHONY: all desktop server build run run-server test install-desktop uninstall-desktop \
        package-arch deb rpm windows macos appimage flatpak snap packages-all \
        systemd-install systemd-uninstall clean help

# Default target: build native Arch Linux desktop application
all: desktop

help:
	@echo "AI-Free YouTube - Universal Packaging & Build Commands:"
	@echo "------------------------------------------------------------------------------"
	@echo "  Arch Linux Desktop:"
	@echo "    make desktop           Build native Wails desktop app (build/bin/aifreeyoutube)"
	@echo "    make run               Launch native desktop app"
	@echo "    make install-desktop   Install desktop app, .desktop entry, and icons to $(PREFIX)"
	@echo "    make uninstall-desktop Remove installed desktop app from $(PREFIX)"
	@echo "    make package-arch      Build official Arch Linux package (.pkg.tar.zst) via makepkg"
	@echo ""
	@echo "  Web Server & Systemd Service:"
	@echo "    make server            Compile standalone Go web server (bin/aifreeyoutube-server)"
	@echo "    make run-server        Launch standalone web server on http://localhost:8484"
	@echo "    make service-start     Start user systemd service (auto-starts on boot)"
	@echo "    make service-stop      Stop user systemd service"
	@echo "    make service-status    Check status of user systemd service"
	@echo "    make systemd-install   Install hardened system-wide systemd service to /etc/systemd/system"
	@echo ""
	@echo "  Cross-Platform Distribution Packages:"
	@echo "    make deb               Build Debian / Ubuntu package (.deb)"
	@echo "    make rpm               Build Red Hat / Fedora / openSUSE package (.rpm)"
	@echo "    make windows           Build Windows 64-bit bundle (.zip / .exe / Chocolatey)"
	@echo "    make macos             Build macOS Apple Silicon & Intel bundles (Homebrew)"
	@echo "    make appimage          Build Linux Universal AppImage"
	@echo "    make flatpak           Prepare/build Flatpak bundle (com.aifreeyoutube.app)"
	@echo "    make snap              Prepare/build Canonical Snap package"
	@echo "    make packages-all      Build all supported distribution formats into build/dist/"
	@echo ""
	@echo "  Quality Assurance:"
	@echo "    make test              Execute comprehensive Go unit test suite"
	@echo "    make clean             Clean compiled binaries and build caches"

# ------------------------------------------------------------------------------
# Arch Linux Desktop Application (Wails v2 + WebKit2GTK 4.1)
# ------------------------------------------------------------------------------

desktop:
	@echo "==> Building Arch Linux native desktop application..."
	cd frontend && npm run build
	$(WAILS) build -tags webkit2_41
	@echo "==> Desktop build complete: build/bin/aifreeyoutube"

run: desktop
	@echo "==> Launching AI-Free YouTube desktop application..."
	./build/bin/aifreeyoutube

install-desktop: desktop
	@echo "==> Installing AI-Free YouTube to $(PREFIX)..."
	install -Dm755 build/bin/aifreeyoutube $(PREFIX)/bin/aifreeyoutube
	install -Dm644 build/linux/aifreeyoutube.desktop $(PREFIX)/share/applications/aifreeyoutube.desktop
	install -Dm644 build/appicon.svg $(PREFIX)/share/icons/hicolor/scalable/apps/aifreeyoutube.svg
	install -Dm644 build/appicon.png $(PREFIX)/share/icons/hicolor/512x512/apps/aifreeyoutube.png
	@which update-desktop-database >/dev/null 2>&1 && update-desktop-database $(PREFIX)/share/applications || true
	@echo "==> Installation complete! You can now launch 'AI-Free YouTube' from your application menu or terminal."

uninstall-desktop:
	@echo "==> Uninstalling AI-Free YouTube from $(PREFIX)..."
	rm -f $(PREFIX)/bin/aifreeyoutube
	rm -f $(PREFIX)/share/applications/aifreeyoutube.desktop
	rm -f $(PREFIX)/share/icons/hicolor/scalable/apps/aifreeyoutube.svg
	rm -f $(PREFIX)/share/icons/hicolor/512x512/apps/aifreeyoutube.png
	@which update-desktop-database >/dev/null 2>&1 && update-desktop-database $(PREFIX)/share/applications || true
	@echo "==> Uninstalled successfully."

package-arch:
	@echo "==> Building Arch Linux package via makepkg..."
	makepkg -f

# ------------------------------------------------------------------------------
# Standalone Web Server & Systemd Services
# ------------------------------------------------------------------------------

server:
	@echo "==> Compiling standalone Go web server..."
	cd frontend && npm run build
	mkdir -p bin
	go build -o bin/aifreeyoutube-server ./cmd/server
	@echo "==> Web server binary compiled: bin/aifreeyoutube-server"

run-server: server
	@echo "==> Launching standalone web server on http://localhost:8484..."
	./bin/aifreeyoutube-server

service-start: server
	@echo "==> Enabling and starting AI-Free YouTube background user service..."
	mkdir -p $(HOME)/.config/systemd/user
	cp packaging/systemd/aifreeyoutube-user.service $(HOME)/.config/systemd/user/aifreeyoutube.service
	systemctl --user daemon-reload
	systemctl --user enable --now aifreeyoutube
	@echo "==> Service active! Web app available at http://localhost:8484"

service-stop:
	@echo "==> Stopping AI-Free YouTube background service..."
	systemctl --user stop aifreeyoutube
	@echo "==> Service stopped."

service-status:
	@systemctl --user status aifreeyoutube

systemd-install: server
	@echo "==> Installing system-wide systemd service (requires sudo)..."
	sudo install -Dm755 bin/aifreeyoutube-server /usr/bin/aifreeyoutube-server
	sudo install -Dm644 packaging/systemd/aifreeyoutube.service /etc/systemd/system/aifreeyoutube.service
	sudo id -u aifreeyoutube >/dev/null 2>&1 || sudo useradd -r -s /bin/false -d /var/lib/aifreeyoutube aifreeyoutube
	sudo mkdir -p /var/lib/aifreeyoutube /var/log/aifreeyoutube
	sudo chown -R aifreeyoutube:aifreeyoutube /var/lib/aifreeyoutube /var/log/aifreeyoutube
	sudo systemctl daemon-reload
	@echo "==> System service installed! Enable with: sudo systemctl enable --now aifreeyoutube"

systemd-uninstall:
	@echo "==> Uninstalling system-wide systemd service..."
	sudo systemctl stop aifreeyoutube || true
	sudo systemctl disable aifreeyoutube || true
	sudo rm -f /etc/systemd/system/aifreeyoutube.service /usr/bin/aifreeyoutube-server
	sudo systemctl daemon-reload
	@echo "==> System service uninstalled."

# ------------------------------------------------------------------------------
# Cross-Platform Packaging Targets
# ------------------------------------------------------------------------------

deb:
	@./packaging/debian/build-deb.sh

rpm:
	@./packaging/rpm/build-rpm.sh

windows:
	@./packaging/windows/build-windows.sh

macos:
	@./packaging/macos/build-macos.sh

appimage:
	@./packaging/appimage/build-appimage.sh

flatpak:
	@./packaging/flatpak/build-flatpak.sh

snap:
	@./packaging/snap/build-snap.sh

packages-all:
	@./packaging/scripts/build-all.sh all

# ------------------------------------------------------------------------------
# Arch User Repository (AUR) Publishing
# ------------------------------------------------------------------------------

aur:
	@echo "==> Updating AUR metadata..."
	@makepkg --printsrcinfo > .SRCINFO
	@mkdir -p packaging/aur
	@cp PKGBUILD .SRCINFO packaging/aur/
	@echo "==> AUR manifests synchronized: PKGBUILD, .SRCINFO"

aur-push: aur
	@./packaging/aur/push-aur.sh

# ------------------------------------------------------------------------------
# Quality Assurance & Cleaning
# ------------------------------------------------------------------------------

test:
	@echo "==> Running test suite..."
	go test -v ./...

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf build/bin
	rm -rf build/dist
	rm -rf build/deb-staging
	rm -rf build/rpm
	rm -rf build/windows-staging*
	rm -rf build/macos-staging*
	rm -rf build/AppDir
	rm -rf bin
	rm -rf frontend/dist
	rm -f *.pkg.tar.zst
	@echo "==> Clean complete."
