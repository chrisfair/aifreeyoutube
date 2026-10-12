package main

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"aifreeyoutube/internal/api"
	"aifreeyoutube/internal/filter"
	"aifreeyoutube/internal/models"
	"aifreeyoutube/internal/youtube"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App manages the lifecycle of the AI-Free YouTube desktop application and provides
// IPC bridges between the Go core engine and the Vue 3 desktop frontend.
//
// Technical Architecture:
// In Wails v2, the Go backend runs as the primary native process, embedding the
// compiled Vue 3 distribution files into WebKit2GTK. This App structure exposes both:
//  1. In-process HTTP routing: The existing REST endpoints are registered with Wails's
//     AssetServer.Handler, allowing all existing frontend fetch calls (`/api/search`, etc.)
//     to execute instantaneously with zero network overhead.
//  2. Direct RPC bindings: Native Go methods bound to `window.go.main.App` in the browser,
//     enabling desktop-native operations such as opening links in the external system browser.
type App struct {
	ctx        context.Context
	ytService  *youtube.Service
	configMgr  *filter.ConfigManager
	detector   *filter.Detector
	apiHandler *api.Handler
	apiRouter  http.Handler
}

// NewApp instantiates the core backend subsystems and configures internal routing.
func NewApp() *App {
	ytService := youtube.NewService()
	configMgr := filter.NewConfigManager()
	detector := filter.NewDetector(configMgr)
	apiHandler := api.NewHandler(ytService, configMgr, detector)
	apiRouter := api.SetupRouter(apiHandler)

	return &App{
		ytService:  ytService,
		configMgr:  configMgr,
		detector:   detector,
		apiHandler: apiHandler,
		apiRouter:  apiRouter,
	}
}

// startup is called when the Wails runtime starts up and creates the native window.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(a.ctx, "AI-Free YouTube desktop engine initialized successfully.")
}

// shutdown is called at application termination.
func (a *App) shutdown(ctx context.Context) {
	runtime.LogInfo(a.ctx, "AI-Free YouTube desktop engine shutting down.")
}

// Search queries YouTube's Innertube API and passes all candidate results through
// the probabilistic AI Content Detector.
func (a *App) Search(query string) (*models.SearchResponse, error) {
	if query == "" {
		return &models.SearchResponse{
			Query:            "",
			Videos:           []models.Video{},
			Quarantined:      []models.Video{},
			HumanFeedRatio:   100.0,
		}, nil
	}

	rawVideos, err := a.ytService.Search(a.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query Innertube search: %w", err)
	}

	settings := a.configMgr.GetSettings()
	clean, quarantined := a.detector.FilterBatch(rawVideos, &settings)
	if clean == nil {
		clean = []models.Video{}
	}
	if quarantined == nil {
		quarantined = []models.Video{}
	}

	total := len(rawVideos)
	ratio := 100.0
	if total > 0 {
		ratio = math.Round((float64(len(clean))/float64(total))*1000) / 10
	}

	return &models.SearchResponse{
		Query:            query,
		TotalScanned:     total,
		TotalClean:       len(clean),
		TotalQuarantined: len(quarantined),
		HumanFeedRatio:   ratio,
		Videos:           clean,
		Quarantined:      quarantined,
	}, nil
}

// GetTrending returns an authentic, human-verified feed of popular videos.
func (a *App) GetTrending() (*models.SearchResponse, error) {
	rawVideos, err := a.ytService.BrowseTrending(a.ctx)
	if err != nil {
		rawVideos, err = a.ytService.BrowseHome(a.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch trending feed: %w", err)
		}
	}

	settings := a.configMgr.GetSettings()
	clean, quarantined := a.detector.FilterBatch(rawVideos, &settings)
	if clean == nil {
		clean = []models.Video{}
	}
	if quarantined == nil {
		quarantined = []models.Video{}
	}

	total := len(rawVideos)
	ratio := 100.0
	if total > 0 {
		ratio = math.Round((float64(len(clean))/float64(total))*1000) / 10
	}

	return &models.SearchResponse{
		Query:            "Trending",
		TotalScanned:     total,
		TotalClean:       len(clean),
		TotalQuarantined: len(quarantined),
		HumanFeedRatio:   ratio,
		Videos:           clean,
		Quarantined:      quarantined,
	}, nil
}

// GetVideo fetches primary metadata and related recommendations for a specific video ID.
func (a *App) GetVideo(id string) (*models.VideoDetailResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("video ID is required")
	}

	primary, related, err := a.ytService.GetVideoDetails(a.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video details: %w", err)
	}

	settings := a.configMgr.GetSettings()
	evaluatedPrimary, _ := a.detector.Evaluate(&primary, &settings)
	cleanRelated, _ := a.detector.FilterBatch(related, &settings)
	if cleanRelated == nil {
		cleanRelated = []models.Video{}
	}

	return &models.VideoDetailResponse{
		Video:         evaluatedPrimary,
		RelatedVideos: cleanRelated,
	}, nil
}

// GetSettings retrieves the active AI shield configuration and sensitivity presets.
func (a *App) GetSettings() models.FilterSettings {
	return a.configMgr.GetSettings()
}

// SaveSettings updates the active filter sensitivity and rule category toggles in memory.
func (a *App) SaveSettings(settings models.FilterSettings) error {
	a.configMgr.UpdateSettings(settings)
	return nil
}

// GetStats returns global lifetime scanning and quarantine metrics.
func (a *App) GetStats() models.GlobalFilterStats {
	return a.detector.GetStats()
}

// BlockChannel adds a channel handle or name to the active blocklist.
func (a *App) BlockChannel(handle string) error {
	a.configMgr.AddBlockedChannel(handle)
	return nil
}

// WhitelistChannel adds a channel handle or name to the trusted human whitelist.
func (a *App) WhitelistChannel(handle string) error {
	a.configMgr.AddWhitelistedChannel(handle)
	return nil
}

// OpenInBrowser launches an external URL (such as a YouTube video link) in Arch Linux's
// default web browser via xdg-open.
func (a *App) OpenInBrowser(url string) error {
	if a.ctx == nil {
		return fmt.Errorf("context not yet initialized")
	}
	runtime.BrowserOpenURL(a.ctx, url)
	return nil
}
