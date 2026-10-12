package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SetupRouter configures the HTTP multiplexer with all API routes, static assets, and middleware.
func SetupRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// Discovery and Content Endpoints
	mux.HandleFunc("/api/search", h.HandleSearch)
	mux.HandleFunc("/api/trending", h.HandleTrending)
	mux.HandleFunc("/api/video", h.HandleVideoDetails)

	// Filter Configuration and Diagnostics Endpoints
	mux.HandleFunc("/api/filter/settings", h.HandleSettings)
	mux.HandleFunc("/api/filter/stats", h.HandleStats)
	mux.HandleFunc("/api/filter/block-channel", h.HandleBlockChannel)
	mux.HandleFunc("/api/filter/whitelist-channel", h.HandleWhitelistChannel)

	// Health Check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "AI-Free YouTube Engine",
		})
	})

	// Static SPA File Server (serves compiled Vue 3 frontend if available)
	distPaths := []string{"../frontend/dist", "frontend/dist", "./dist"}
	var distDir string
	for _, p := range distPaths {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			distDir = p
			break
		}
	}

	if distDir != "" {
		fs := http.FileServer(http.Dir(distDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}
			// If a static file exists (e.g. .css, .js, .ico), serve it directly
			targetFile := filepath.Join(distDir, filepath.Clean(r.URL.Path))
			if fi, err := os.Stat(targetFile); err == nil && !fi.IsDir() {
				fs.ServeHTTP(w, r)
				return
			}
			// Otherwise fallback to index.html to support Vue Router HTML5 history mode
			http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
		})
	}

	// Wrap mux with CORS and request logging middleware
	return EnableCORS(RequestLogger(mux))
}
