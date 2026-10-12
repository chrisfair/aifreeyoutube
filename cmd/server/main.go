// Package main starts the AI-Free YouTube backend service.
// It initializes the YouTube Innertube client, mounts the multi-signal AI filter engine,
// configures REST endpoints, and manages graceful HTTP server lifecycle.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aifreeyoutube/internal/api"
	"aifreeyoutube/internal/filter"
	"aifreeyoutube/internal/youtube"
)

func main() {
	log.Println("==================================================")
	log.Println("       AI-Free YouTube - Core Engine v1.0.0       ")
	log.Println("==================================================")
	log.Println("Initializing YouTube Innertube client...")
	ytService := youtube.NewService()

	log.Println("Initializing AI Content Detection Engine...")
	configMgr := filter.NewConfigManager()
	detector := filter.NewDetector(configMgr)

	log.Println("Configuring API handlers and middleware...")
	handler := api.NewHandler(ytService, configMgr, detector)
	router := api.SetupRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background goroutine
	go func() {
		log.Printf("Server listening on http://localhost:%s", port)
		log.Printf("Endpoints:")
		log.Printf("  • GET  http://localhost:%s/api/trending", port)
		log.Printf("  • GET  http://localhost:%s/api/search?q={query}", port)
		log.Printf("  • GET  http://localhost:%s/api/video?id={videoId}", port)
		log.Printf("  • GET  http://localhost:%s/api/filter/settings", port)
		log.Printf("  • GET  http://localhost:%s/api/filter/stats", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	// Listen for interrupt signals to enable graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down AI-Free YouTube server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
