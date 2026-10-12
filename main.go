// Package main serves as the primary entrypoint for the AI-Free YouTube Arch Linux desktop application.
// It bundles the Vue 3 frontend single-page application into native WebKit2GTK windows via Wails v2.
package main

import (
	"embed"
	_ "embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

func init() {
	// Wayland + NVIDIA / Hyprland Compatibility:
	// On Wayland compositors (especially Hyprland, Sway, and GNOME Wayland with NVIDIA drivers),
	// WebKit2GTK's DMA-BUF renderer triggers "Error 71 (Protocol error) dispatching to Wayland display".
	// Setting WEBKIT_DISABLE_DMABUF_RENDERER=1 disables the problematic DMA-BUF subsurface pathway,
	// allowing WebKit2GTK to render flawlessly on Wayland without protocol crashes.
	if os.Getenv("WEBKIT_DISABLE_DMABUF_RENDERER") == "" {
		_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
	}
}

// assets embeds the compiled Vue 3 production distribution.
//
//go:embed all:frontend/dist
var assets embed.FS

// appIcon embeds the 512x512 PNG application icon for native Linux window managers.
//
//go:embed build/appicon.png
var appIcon []byte

func main() {
	// Create an instance of the backend application controller
	app := NewApp()

	// Launch the native Wails desktop application
	err := wails.Run(&options.App{
		Title:     "AI-Free YouTube - Authentic Human Content Portal",
		Width:     1366,
		Height:    860,
		MinWidth:  960,
		MinHeight: 640,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: app.apiRouter, // In-process routing for all /api/* requests
		},
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 1}, // Dark obsidian bg-surface-950
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			Icon:                appIcon,
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways, // Hardware acceleration on Linux
			ProgramName:         "aifreeyoutube",
		},
	})

	if err != nil {
		log.Fatalf("Fatal error running AI-Free YouTube desktop app: %v", err)
	}
}
