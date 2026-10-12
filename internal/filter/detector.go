package filter

import (
	"math"
	"sync"
	"time"

	"aifreeyoutube/internal/models"
)

// Detector is the central evaluation engine that runs incoming videos through
// lexical heuristics, clickbait templates, channel identifiers, and blocklists.
type Detector struct {
	configMgr    *ConfigManager
	blocklistMgr *BlocklistManager

	statsMu                sync.RWMutex
	totalVideosScanned     int64
	totalVideosQuarantined int64
	totalVideosPassed      int64
	startTime              time.Time
}

// NewDetector initializes the AI Content Detection Engine.
func NewDetector(configMgr *ConfigManager) *Detector {
	return &Detector{
		configMgr:    configMgr,
		blocklistMgr: NewBlocklistManager(),
		startTime:    time.Now(),
	}
}

// Evaluate performs a comprehensive multi-signal inspection of a single video item.
// It returns the enriched video (with AIScore, HumanScore, Verdict, and Signals)
// and a boolean indicating whether the video passed (true = clean, false = quarantined).
func (d *Detector) Evaluate(video *models.Video, settings *models.FilterSettings) (models.Video, bool) {
	d.recordScan()

	evaluated := *video
	var signals []models.FilterSignal

	// 1. Channel Evaluation & Whitelist Check
	chanSignals, isWhitelisted := EvaluateChannel(
		video.ChannelName,
		video.ChannelID,
		video.ChannelHandle,
		d.blocklistMgr,
		settings,
	)

	// If creator is whitelisted by user, bypass all AI checks
	if isWhitelisted {
		evaluated.AIScore = 0.0
		evaluated.HumanScore = 1.0
		evaluated.Verdict = models.VerdictClean
		evaluated.Signals = []models.FilterSignal{
			{
				RuleName:    "Whitelisted Creator",
				Category:    "user_rule",
				Weight:      0.0,
				Confidence:  1.0,
				Description: "Channel is explicitly verified and whitelisted by user",
			},
		}
		d.recordPass()
		return evaluated, true
	}

	// Append channel signals if enabled
	if settings.EnabledCategories["channel"] || settings.EnabledCategories["blocklist"] {
		signals = append(signals, chanSignals...)
	}

	// 2. Custom User Blacklist Keywords
	customSignals := EvaluateCustomKeywords(video.Title, video.Description, settings.CustomBlockedKeywords)
	signals = append(signals, customSignals...)

	// 3. YouTube Shorts Filter
	if settings.FilterShorts && video.IsShort {
		signals = append(signals, models.FilterSignal{
			RuleName:    "YouTube Shorts Format",
			Category:    "metadata",
			Weight:      0.85,
			Confidence:  0.95,
			Description: "Shorts format eliminated by active user filter preference",
			MatchedText: "#shorts",
		})
	}

	// 4. Text & Lexical Heuristics (Title & Description)
	textSignals := EvaluateText(video.Title, video.Description)
	for _, sig := range textSignals {
		if settings.EnabledCategories[sig.Category] {
			signals = append(signals, sig)
		}
	}

	// 5. Aggregate Signals into Combined AI Probability Score
	// We use probabilistic OR accumulation: combined = 1 - Π(1 - (weight * confidence))
	// This ensures multiple signals compound smoothly without exceeding 1.0.
	accumulatedUncertainty := 1.0
	for _, s := range signals {
		factor := s.Weight * s.Confidence
		if factor > 0.99 {
			factor = 0.99
		}
		accumulatedUncertainty *= (1.0 - factor)
	}

	aiScore := 1.0 - accumulatedUncertainty
	// Round to two decimal places
	aiScore = math.Round(aiScore*100) / 100

	evaluated.AIScore = aiScore
	evaluated.HumanScore = math.Round((1.0-aiScore)*100) / 100
	evaluated.Signals = signals

	// 6. Assign Categorical Verdict based on computed score and active threshold
	if evaluated.AIScore >= settings.Threshold {
		evaluated.Verdict = models.VerdictAISlop
		d.recordQuarantine()
		return evaluated, false // Quarantined
	} else if evaluated.AIScore >= 0.35 {
		evaluated.Verdict = models.VerdictSuspicious
		d.recordPass()
		return evaluated, true // Allowed with caution badge
	} else {
		evaluated.Verdict = models.VerdictClean
		d.recordPass()
		return evaluated, true // Clean
	}
}

// FilterBatch processes a slice of raw videos, separating them into Clean and Quarantined arrays.
func (d *Detector) FilterBatch(videos []models.Video, settings *models.FilterSettings) ([]models.Video, []models.Video) {
	var clean []models.Video
	var quarantined []models.Video

	for _, v := range videos {
		enriched, passed := d.Evaluate(&v, settings)
		if passed {
			clean = append(clean, enriched)
		} else {
			quarantined = append(quarantined, enriched)
		}
	}

	return clean, quarantined
}

// GetStats returns global runtime metrics across all requests.
func (d *Detector) GetStats() models.GlobalFilterStats {
	d.statsMu.RLock()
	defer d.statsMu.RUnlock()

	return models.GlobalFilterStats{
		TotalVideosScanned:     d.totalVideosScanned,
		TotalVideosQuarantined: d.totalVideosQuarantined,
		TotalVideosPassed:      d.totalVideosPassed,
		ServerUptimeSeconds:    int64(time.Since(d.startTime).Seconds()),
		LastUpdated:            time.Now(),
	}
}

func (d *Detector) recordScan() {
	d.statsMu.Lock()
	defer d.statsMu.Unlock()
	d.totalVideosScanned++
}

func (d *Detector) recordPass() {
	d.statsMu.Lock()
	defer d.statsMu.Unlock()
	d.totalVideosPassed++
}

func (d *Detector) recordQuarantine() {
	d.statsMu.Lock()
	defer d.statsMu.Unlock()
	d.totalVideosQuarantined++
}
