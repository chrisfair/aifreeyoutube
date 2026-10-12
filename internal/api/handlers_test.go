package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"aifreeyoutube/internal/filter"
	"aifreeyoutube/internal/models"
	"aifreeyoutube/internal/youtube"
)

func setupTestRouter() http.Handler {
	ytService := youtube.NewService()
	configMgr := filter.NewConfigManager()
	detector := filter.NewDetector(configMgr)
	handler := NewHandler(ytService, configMgr, detector)
	return SetupRouter(handler)
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", resp["status"])
	}
}

func TestSettingsEndpoints(t *testing.T) {
	router := setupTestRouter()

	// 1. GET Settings
	getReq := httptest.NewRequest(http.MethodGet, "/api/filter/settings", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /api/filter/settings returned status %d", getRec.Code)
	}

	var settings models.FilterSettings
	if err := json.NewDecoder(getRec.Body).Decode(&settings); err != nil {
		t.Fatalf("failed to decode settings: %v", err)
	}

	if settings.SensitivityMode != "balanced" {
		t.Errorf("expected default sensitivity 'balanced', got '%s'", settings.SensitivityMode)
	}

	// 2. POST Settings (Update to Aggressive)
	settings.SensitivityMode = "aggressive"
	payload, _ := json.Marshal(settings)

	postReq := httptest.NewRequest(http.MethodPost, "/api/filter/settings", bytes.NewBuffer(payload))
	postReq.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()
	router.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("POST /api/filter/settings returned status %d", postRec.Code)
	}

	var updatedSettings models.FilterSettings
	if err := json.NewDecoder(postRec.Body).Decode(&updatedSettings); err != nil {
		t.Fatalf("failed to decode updated settings: %v", err)
	}

	if updatedSettings.SensitivityMode != "aggressive" || updatedSettings.Threshold != filter.PresetAggressiveThreshold {
		t.Errorf("expected aggressive threshold %f, got %f", filter.PresetAggressiveThreshold, updatedSettings.Threshold)
	}
}

func TestBlockAndWhitelistEndpoints(t *testing.T) {
	router := setupTestRouter()

	// 1. Block Channel
	blockPayload, _ := json.Marshal(map[string]string{"channel": "UC_spammer_bot"})
	blockReq := httptest.NewRequest(http.MethodPost, "/api/filter/block-channel", bytes.NewBuffer(blockPayload))
	blockReq.Header.Set("Content-Type", "application/json")
	blockRec := httptest.NewRecorder()
	router.ServeHTTP(blockRec, blockReq)

	if blockRec.Code != http.StatusOK {
		t.Fatalf("POST /api/filter/block-channel returned status %d", blockRec.Code)
	}

	// 2. Whitelist Channel
	whitePayload, _ := json.Marshal(map[string]string{"channel": "UC_awesome_creator"})
	whiteReq := httptest.NewRequest(http.MethodPost, "/api/filter/whitelist-channel", bytes.NewBuffer(whitePayload))
	whiteReq.Header.Set("Content-Type", "application/json")
	whiteRec := httptest.NewRecorder()
	router.ServeHTTP(whiteRec, whiteReq)

	if whiteRec.Code != http.StatusOK {
		t.Fatalf("POST /api/filter/whitelist-channel returned status %d", whiteRec.Code)
	}
}
