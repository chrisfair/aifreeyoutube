// Package models defines data structures used across the AI-Free YouTube application.
// This includes video representations, Innertube extraction schemas, AI scoring metrics,
// and user configuration profiles for filtering synthetic and robotic content.
package models

import "time"

// FilterVerdict defines the classification assigned to a video after evaluation
// by the AI Content Detection Engine.
type FilterVerdict string

const (
	// VerdictClean indicates the video passed inspection with minimal to no synthetic indicators.
	VerdictClean FilterVerdict = "CLEAN"

	// VerdictSuspicious indicates the video has multiple borderline markers that warrant user caution.
	VerdictSuspicious FilterVerdict = "SUSPICIOUS"

	// VerdictAISlop indicates the video exceeds the threshold for synthetic/robotic content generation.
	VerdictAISlop FilterVerdict = "AI_SLOP"
)

// FilterSignal represents a single heuristic rule trigger that contributed to a video's AI score.
// Each signal provides transparent, explainable diagnostics to the user regarding why content was flagged.
type FilterSignal struct {
	// RuleName is the unique human-readable identifier of the heuristic rule (e.g., "ChatGPT Buzzword Pattern").
	RuleName string `json:"ruleName"`

	// Category groups the rule type: "lexical", "formula", "channel", "metadata", or "blocklist".
	Category string `json:"category"`

	// Weight is the contribution factor of this rule towards the overall AI probability score (0.0 to 1.0).
	Weight float64 `json:"weight"`

	// Confidence indicates the mathematical certainty of the pattern match (0.0 to 1.0).
	Confidence float64 `json:"confidence"`

	// Description provides human-readable context explaining why the trigger occurred.
	Description string `json:"description"`

	// MatchedText contains the specific keyword, phrase, or token that caused the trigger.
	MatchedText string `json:"matchedText,omitempty"`
}

// Video represents an enriched YouTube video item processed by our backend.
// It bundles standard metadata extracted via Innertube alongside proprietary AI detection metrics.
type Video struct {
	// ID is the YouTube 11-character unique video identifier.
	ID string `json:"id"`

	// Title is the video's headline.
	Title string `json:"title"`

	// ChannelName is the display name of the publishing channel or creator.
	ChannelName string `json:"channelName"`

	// ChannelID is the YouTube channel browse ID (e.g., "UC...").
	ChannelID string `json:"channelId"`

	// ChannelHandle is the @handle of the channel if available.
	ChannelHandle string `json:"channelHandle,omitempty"`

	// Views represents the formatted view count string (e.g., "1.2M views").
	Views string `json:"views"`

	// Published is the relative upload timestamp (e.g., "3 days ago").
	Published string `json:"published"`

	// Duration represents the video run time (e.g., "14:22").
	Duration string `json:"duration"`

	// Description is the excerpt or full summary provided with the video.
	Description string `json:"description"`

	// Thumbnail is the high-resolution image URL for the video preview.
	Thumbnail string `json:"thumbnail"`

	// EmbedURL is the privacy-friendly YouTube nocookie embed link used by the player.
	EmbedURL string `json:"embedUrl"`

	// IsShort indicates whether the video is a YouTube Short format (< 60s vertical).
	IsShort bool `json:"isShort"`

	// AIScore represents the computed probability that the content is primarily synthetic/AI-generated (0.0 to 1.0).
	// 0.0 indicates 100% human authenticity, while 1.0 indicates definite synthetic slop.
	AIScore float64 `json:"aiScore"`

	// HumanScore is the inverse of AIScore (1.0 - AIScore), providing a positive metric for clean feeds.
	HumanScore float64 `json:"humanScore"`

	// Verdict is the categorized status ("CLEAN", "SUSPICIOUS", "AI_SLOP").
	Verdict FilterVerdict `json:"verdict"`

	// Signals contains the detailed list of heuristic rules triggered during evaluation.
	Signals []FilterSignal `json:"signals"`
}

// FilterSettings stores customizable user preferences that govern how aggressively
// synthetic content is filtered across the application.
type FilterSettings struct {
	// SensitivityMode defines presets: "aggressive" (0.35 threshold), "balanced" (0.55), or "mild" (0.75).
	SensitivityMode string `json:"sensitivityMode"`

	// Threshold is the custom cutoff score above which videos are quarantined (0.0 to 1.0). Default is 0.55.
	Threshold float64 `json:"threshold"`

	// FilterShorts toggles whether YouTube Shorts format videos are automatically eliminated.
	FilterShorts bool `json:"filterShorts"`

	// EnabledCategories toggles individual rule categories (e.g., "lexical": true, "channel": true).
	EnabledCategories map[string]bool `json:"enabledCategories"`

	// CustomBlockedKeywords is a user-managed list of words/phrases that immediately flag a video.
	CustomBlockedKeywords []string `json:"customBlockedKeywords"`

	// CustomBlockedChannels is a user-managed list of channel IDs or handles to exclude.
	CustomBlockedChannels []string `json:"customBlockedChannels"`

	// WhitelistedChannels is a user-managed list of channel IDs exempt from AI filtering.
	WhitelistedChannels []string `json:"whitelistedChannels"`
}

// SearchResponse models the JSON payload returned to the frontend upon searching or browsing.
type SearchResponse struct {
	// Query is the search term used, or "Trending" for discovery feeds.
	Query string `json:"query"`

	// TotalScanned is the total number of raw videos evaluated from YouTube.
	TotalScanned int `json:"totalScanned"`

	// TotalClean is the count of videos that successfully passed the AI filter.
	TotalClean int `json:"totalClean"`

	// TotalQuarantined is the count of synthetic/AI videos intercepted and blocked.
	TotalQuarantined int `json:"totalQuarantined"`

	// HumanFeedRatio is the percentage of evaluated content that was clean (TotalClean / TotalScanned).
	HumanFeedRatio float64 `json:"humanFeedRatio"`

	// Videos is the array of clean, verified videos intended for direct display in the feed.
	Videos []Video `json:"videos"`

	// Quarantined contains videos that triggered the AI filter, available in the Quarantine Drawer for transparency.
	Quarantined []Video `json:"quarantined"`
}

// VideoDetailResponse models the complete metadata returned when watching a specific video.
type VideoDetailResponse struct {
	// Video contains the core video details.
	Video Video `json:"video"`

	// RelatedVideos contains clean, AI-filtered recommendations for the watch page sidebar.
	RelatedVideos []Video `json:"relatedVideos"`

	// QuarantinedRelated contains recommendations that were blocked from the sidebar due to synthetic content.
	QuarantinedRelated []Video `json:"quarantinedRelated"`
}

// GlobalFilterStats provides aggregate lifetime metrics across all search and browsing sessions.
type GlobalFilterStats struct {
	// TotalVideosScanned is the lifetime count of videos processed through the filter engine.
	TotalVideosScanned int64 `json:"totalVideosScanned"`

	// TotalVideosQuarantined is the lifetime count of AI/synthetic videos blocked.
	TotalVideosQuarantined int64 `json:"totalVideosQuarantined"`

	// TotalVideosPassed is the lifetime count of human-verified videos delivered.
	TotalVideosPassed int64 `json:"totalVideosPassed"`

	// ServerUptimeSeconds tracks how long the backend service has been running.
	ServerUptimeSeconds int64 `json:"serverUptimeSeconds"`

	// LastUpdated is the timestamp of the last statistic update.
	LastUpdated time.Time `json:"lastUpdated"`
}
