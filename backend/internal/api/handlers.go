package api

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"

	"aifreeyoutube/internal/filter"
	"aifreeyoutube/internal/models"
	"aifreeyoutube/internal/youtube"
)

// Handler manages HTTP API requests for search, feeds, video inspection, and filter settings.
type Handler struct {
	ytService *youtube.Service
	configMgr *filter.ConfigManager
	detector  *filter.Detector
}

// NewHandler initializes the API handler dependencies.
func NewHandler(ytService *youtube.Service, configMgr *filter.ConfigManager, detector *filter.Detector) *Handler {
	return &Handler{
		ytService: ytService,
		configMgr: configMgr,
		detector:  detector,
	}
}

// respondJSON is a utility that writes a JSON response with the appropriate content-type and HTTP status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// respondError formats an error message into a standardized JSON payload.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// resolveSettings blends global server filter settings with optional per-request sensitivity overrides.
func (h *Handler) resolveSettings(r *http.Request) models.FilterSettings {
	settings := h.configMgr.GetSettings()
	overrideSensitivity := strings.ToLower(r.URL.Query().Get("sensitivity"))

	switch overrideSensitivity {
	case "aggressive":
		settings.SensitivityMode = "aggressive"
		settings.Threshold = filter.PresetAggressiveThreshold
	case "balanced":
		settings.SensitivityMode = "balanced"
		settings.Threshold = filter.PresetBalancedThreshold
	case "mild":
		settings.SensitivityMode = "mild"
		settings.Threshold = filter.PresetMildThreshold
	}

	return settings
}

// HandleSearch executes a YouTube query, filters the resulting batch with the AI detector,
// and separates clean videos from quarantined slop.
// GET /api/search?q={query}&sensitivity={level}
func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		respondError(w, http.StatusBadRequest, "Missing required 'q' search parameter")
		return
	}

	settings := h.resolveSettings(r)

	rawVideos, err := h.ytService.Search(r.Context(), query)
	if err != nil {
		log.Printf("YouTube search error for '%s': %v", query, err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve search results from YouTube")
		return
	}

	clean, quarantined := h.detector.FilterBatch(rawVideos, &settings)
	if clean == nil {
		clean = []models.Video{}
	}
	if quarantined == nil {
		quarantined = []models.Video{}
	}

	totalScanned := len(rawVideos)
	var humanRatio float64
	if totalScanned > 0 {
		humanRatio = math.Round((float64(len(clean))/float64(totalScanned))*1000) / 10
	}

	resp := models.SearchResponse{
		Query:            query,
		TotalScanned:     totalScanned,
		TotalClean:       len(clean),
		TotalQuarantined: len(quarantined),
		HumanFeedRatio:   humanRatio,
		Videos:           clean,
		Quarantined:      quarantined,
	}

	respondJSON(w, http.StatusOK, resp)
}

// HandleTrending retrieves YouTube's trending / explore feed and passes it through the AI filter.
// GET /api/trending?sensitivity={level}
func (h *Handler) HandleTrending(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	settings := h.resolveSettings(r)

	rawVideos, err := h.ytService.BrowseTrending(r.Context())
	if err != nil {
		log.Printf("Trending feed error: %v, attempting home feed fallback...", err)
		// Fallback to home feed if trending is empty or blocked
		rawVideos, err = h.ytService.BrowseHome(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to retrieve discovery feed from YouTube")
			return
		}
	}

	clean, quarantined := h.detector.FilterBatch(rawVideos, &settings)
	if clean == nil {
		clean = []models.Video{}
	}
	if quarantined == nil {
		quarantined = []models.Video{}
	}

	totalScanned := len(rawVideos)
	var humanRatio float64
	if totalScanned > 0 {
		humanRatio = math.Round((float64(len(clean))/float64(totalScanned))*1000) / 10
	}

	resp := models.SearchResponse{
		Query:            "Trending",
		TotalScanned:     totalScanned,
		TotalClean:       len(clean),
		TotalQuarantined: len(quarantined),
		HumanFeedRatio:   humanRatio,
		Videos:           clean,
		Quarantined:      quarantined,
	}

	respondJSON(w, http.StatusOK, resp)
}

// HandleVideoDetails fetches full metadata for a video and cleans its recommendation sidebar.
// GET /api/video?id={id}&sensitivity={level}
func (h *Handler) HandleVideoDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	videoID := strings.TrimSpace(r.URL.Query().Get("id"))
	if videoID == "" {
		respondError(w, http.StatusBadRequest, "Missing required 'id' parameter")
		return
	}

	settings := h.resolveSettings(r)

	primaryVideo, relatedVideos, err := h.ytService.GetVideoDetails(r.Context(), videoID)
	if err != nil {
		log.Printf("Video details error for '%s': %v", videoID, err)
		respondError(w, http.StatusInternalServerError, "Failed to retrieve video details from YouTube")
		return
	}

	// Evaluate the primary video itself so the user sees its AI Shield breakdown on the watch page
	evaluatedPrimary, _ := h.detector.Evaluate(&primaryVideo, &settings)

	cleanRelated, quarantinedRelated := h.detector.FilterBatch(relatedVideos, &settings)
	if cleanRelated == nil {
		cleanRelated = []models.Video{}
	}
	if quarantinedRelated == nil {
		quarantinedRelated = []models.Video{}
	}

	resp := models.VideoDetailResponse{
		Video:              evaluatedPrimary,
		RelatedVideos:      cleanRelated,
		QuarantinedRelated: quarantinedRelated,
	}

	respondJSON(w, http.StatusOK, resp)
}

// HandleSettings retrieves or updates the user's filter preferences.
// GET /api/filter/settings or POST /api/filter/settings
func (h *Handler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		settings := h.configMgr.GetSettings()
		respondJSON(w, http.StatusOK, settings)

	case http.MethodPost:
		var newSettings models.FilterSettings
		if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON settings payload")
			return
		}
		updated := h.configMgr.UpdateSettings(newSettings)
		respondJSON(w, http.StatusOK, updated)

	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleStats returns lifetime scan and quarantine metrics.
// GET /api/filter/stats
func (h *Handler) HandleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	stats := h.detector.GetStats()
	respondJSON(w, http.StatusOK, stats)
}

// HandleBlockChannel adds a channel to the custom blocklist.
// POST /api/filter/block-channel
func (h *Handler) HandleBlockChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Channel) == "" {
		respondError(w, http.StatusBadRequest, "Channel identifier required")
		return
	}

	h.configMgr.AddBlockedChannel(req.Channel)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Channel added to blocklist",
		"channel": req.Channel,
	})
}

// HandleWhitelistChannel adds a creator to the user's whitelist.
// POST /api/filter/whitelist-channel
func (h *Handler) HandleWhitelistChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Channel) == "" {
		respondError(w, http.StatusBadRequest, "Channel identifier required")
		return
	}

	h.configMgr.AddWhitelistedChannel(req.Channel)
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Channel added to whitelist",
		"channel": req.Channel,
	})
}
