// Package filter implements the multi-signal AI Content Detection Engine for YouTube.
// It evaluates titles, descriptions, channel characteristics, tags, and metadata to identify
// synthetic, robotic, and algorithmically generated content mills (slop).
package filter

import (
	"sync"

	"aifreeyoutube/internal/models"
)

const (
	// PresetAggressiveThreshold is the strict cutoff (35% AI score).
	// Suitable for users who want zero tolerance for synthetic scripts, stock voiceovers, or clickbait.
	PresetAggressiveThreshold = 0.35

	// PresetBalancedThreshold is the default cutoff (55% AI score).
	// Accurately intercepts obvious AI content farms while preserving authentic human creators discussing technology.
	PresetBalancedThreshold = 0.55

	// PresetMildThreshold is the conservative cutoff (75% AI score).
	// Only blocks blatant, automated spam mills with multiple high-confidence triggers.
	PresetMildThreshold = 0.75
)

// ConfigManager maintains runtime filter preferences in a thread-safe manner,
// allowing users to dynamically adjust sensitivity, toggle rules, and manage whitelists/blocklists.
type ConfigManager struct {
	mu       sync.RWMutex
	settings models.FilterSettings
}

// NewConfigManager initializes default filter configuration with balanced sensitivity.
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		settings: models.FilterSettings{
			SensitivityMode: "balanced",
			Threshold:       PresetBalancedThreshold,
			FilterShorts:    true, // Shorts are overwhelmingly dominated by AI brainrot/TTS narration
			EnabledCategories: map[string]bool{
				"lexical":   true, // LLM vocabulary and phrasing checks
				"formula":   true, // Cliché automated clickbait template checks
				"metadata":  true, // AI generation tags and hashtag spam
				"channel":   true, // Suspicious channel name/handle patterns
				"blocklist": true, // Preloaded known AI farm database
			},
			CustomBlockedKeywords: []string{},
			CustomBlockedChannels: []string{},
			WhitelistedChannels:   []string{},
		},
	}
}

// GetSettings returns a snapshot of the current filter settings.
func (cm *ConfigManager) GetSettings() models.FilterSettings {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Return a deep copy to prevent race conditions
	copied := cm.settings
	copied.EnabledCategories = make(map[string]bool, len(cm.settings.EnabledCategories))
	for k, v := range cm.settings.EnabledCategories {
		copied.EnabledCategories[k] = v
	}
	copied.CustomBlockedKeywords = append([]string(nil), cm.settings.CustomBlockedKeywords...)
	copied.CustomBlockedChannels = append([]string(nil), cm.settings.CustomBlockedChannels...)
	copied.WhitelistedChannels = append([]string(nil), cm.settings.WhitelistedChannels...)

	return copied
}

// UpdateSettings updates user preferences, validating thresholds and updating sensitivity modes.
func (cm *ConfigManager) UpdateSettings(newSettings models.FilterSettings) models.FilterSettings {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Validate and synchronize threshold with sensitivity mode preset if specified
	switch newSettings.SensitivityMode {
	case "aggressive":
		newSettings.Threshold = PresetAggressiveThreshold
	case "balanced":
		newSettings.Threshold = PresetBalancedThreshold
	case "mild":
		newSettings.Threshold = PresetMildThreshold
	default:
		// Custom threshold clamped between 0.05 and 0.95
		if newSettings.Threshold < 0.05 {
			newSettings.Threshold = 0.05
		}
		if newSettings.Threshold > 0.95 {
			newSettings.Threshold = 0.95
		}
		newSettings.SensitivityMode = "custom"
	}

	if newSettings.EnabledCategories == nil {
		newSettings.EnabledCategories = cm.settings.EnabledCategories
	}

	cm.settings = newSettings
	return cm.settings
}

// AddBlockedChannel appends a channel ID or handle to the user's custom blocklist.
func (cm *ConfigManager) AddBlockedChannel(channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, ch := range cm.settings.CustomBlockedChannels {
		if ch == channel {
			return
		}
	}
	cm.settings.CustomBlockedChannels = append(cm.settings.CustomBlockedChannels, channel)
}

// AddWhitelistedChannel exempts a channel from AI filtering.
func (cm *ConfigManager) AddWhitelistedChannel(channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, ch := range cm.settings.WhitelistedChannels {
		if ch == channel {
			return
		}
	}
	cm.settings.WhitelistedChannels = append(cm.settings.WhitelistedChannels, channel)
}
