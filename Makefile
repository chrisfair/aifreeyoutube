# ==============================================================================
# AI-Free YouTube: Arch Linux Desktop & Web Application
# ==============================================================================
# Makefile for building, packaging, running, and installing the application.
# ==============================================================================

WAILS := $(shell which wails 2>/dev/null || echo $(HOME)/go/bin/wails)
PREFIX ?= $(HOME)/.local

.PHONY: all desktop server build run run-server test install-desktop uninstall-desktop package-arch clean help

# Default target: build native Arch Linux desktop application
all: desktop

help:
	@echo "AI-Free YouTube - Available Commands:"
	@echo "  make desktop           Build native Arch Linux Wails desktop app (build/bin/aifreeyoutube)"
	@echo "  make run               Launch the native Wails desktop application"
	@echo "  make install-desktop   Install desktop app, .desktop entry, and icons to $(PREFIX)"
	@echo "  make uninstall-desktop Remove installed desktop app from $(PREFIX)"
	@echo "  make package-arch      Build official Arch Linux package (.pkg.tar.zst) via makepkg"
	@echo "  make server            Compile standalone Go web server (bin/aifreeyoutube-server)"
	@echo "  make run-server        Launch standalone web server on http://localhost:8080"
	@echo "  make test              Execute comprehensive Go unit test suite"
	@echo "  make clean             Clean compiled binaries and build caches"

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
# Standalone Web Server
# ------------------------------------------------------------------------------

server:
	@echo "==> Compiling standalone Go web server..."
	cd frontend && npm run build
	mkdir -p bin
	go build -o bin/aifreeyoutube-server ./cmd/server
	@echo "==> Web server binary compiled: bin/aifreeyoutube-server"

run-server: server
	@echo "==> Launching standalone web server on http://localhost:8080..."
	./bin/aifreeyoutube-server

service-start: server
	@echo "==> Enabling and starting AI-Free YouTube background service..."
	systemctl --user daemon-reload
	systemctl --user enable --now aifreeyoutube
	@echo "==> Service active! Web app available at http://localhost:8080"

service-stop:
	@echo "==> Stopping AI-Free YouTube background service..."
	systemctl --user stop aifreeyoutube
	@echo "==> Service stopped."

service-status:
	@systemctl --user status aifreeyoutube

# ------------------------------------------------------------------------------
# Quality Assurance & Cleaning
# ------------------------------------------------------------------------------

test:
	@echo "==> Running test suite..."
	go test -v ./...

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf build/bin
	rm -rf bin
	rm -rf frontend/dist
	rm -f *.pkg.tar.zst
	@echo "==> Clean complete."
