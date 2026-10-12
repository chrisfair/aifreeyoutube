package filter

import (
	"regexp"
	"strings"

	"aifreeyoutube/internal/models"
)

var (
	// SuspiciousHandleRegex matches bot handle patterns like "@HistoryFactsAI", "@StoicWisdom_ai", etc.
	SuspiciousHandleRegex = regexp.MustCompile(`(?i)(_ai$|ai_official$|fact[sz]_ai|quotes_ai|recaps?_ai|automation_ai)`)

	// GeneratedBotHandleRegex matches standard YouTube auto-generated bot handles like "@user-ab3df7kz"
	GeneratedBotHandleRegex = regexp.MustCompile(`(?i)^user-[a-z0-9]{8,12}$`)
)

// EvaluateChannel inspects channel identity, handle patterns, and blocklist memberships.
func EvaluateChannel(
	channelName, channelID, channelHandle string,
	blocklistMgr *BlocklistManager,
	settings *models.FilterSettings,
) ([]models.FilterSignal, bool) {
	var signals []models.FilterSignal

	// Check if channel is explicitly whitelisted by user
	if blocklistMgr.IsWhitelisted(channelID, channelHandle, settings.WhitelistedChannels) {
		return nil, true // Whitelist override active
	}

	// 1. Check known AI farm blocklist
	if isFarm, reason := blocklistMgr.MatchKnownFarm(channelName, channelHandle); isFarm {
		signals = append(signals, models.FilterSignal{
			RuleName:    "Known AI Content Farm",
			Category:    "blocklist",
			Weight:      0.80,
			Confidence:  0.95,
			Description: reason,
			MatchedText: channelName,
		})
	}

	// 2. Check user's custom blocked channels
	if blocklistMgr.IsCustomBlocked(channelName, channelID, channelHandle, settings.CustomBlockedChannels) {
		signals = append(signals, models.FilterSignal{
			RuleName:    "Custom Blocked Channel",
			Category:    "blocklist",
			Weight:      1.0,
			Confidence:  1.0,
			Description: "Channel matches user's custom blocklist entry",
			MatchedText: channelName,
		})
	}

	// 3. Inspect channel handle heuristics
	if SuspiciousHandleRegex.MatchString(channelHandle) {
		signals = append(signals, models.FilterSignal{
			RuleName:    "Synthetic Channel Handle",
			Category:    "channel",
			Weight:      0.45,
			Confidence:  0.85,
			Description: "Channel handle explicitly denotes automated AI generation: " + channelHandle,
			MatchedText: channelHandle,
		})
	} else if GeneratedBotHandleRegex.MatchString(channelHandle) {
		signals = append(signals, models.FilterSignal{
			RuleName:    "Auto-Generated Bot Handle",
			Category:    "channel",
			Weight:      0.30,
			Confidence:  0.70,
			Description: "Channel utilizes an unconfigured, auto-generated bot handle identifier: " + channelHandle,
			MatchedText: channelHandle,
		})
	}

	return signals, false
}

// EvaluateCustomKeywords checks if any user-defined keywords appear in title or description.
func EvaluateCustomKeywords(title, description string, customKeywords []string) []models.FilterSignal {
	var signals []models.FilterSignal
	combined := strings.ToLower(title + " " + description)

	for _, kw := range customKeywords {
		kwTrimmed := strings.TrimSpace(strings.ToLower(kw))
		if kwTrimmed == "" {
			continue
		}
		if strings.Contains(combined, kwTrimmed) {
			signals = append(signals, models.FilterSignal{
				RuleName:    "Custom Blocked Keyword",
				Category:    "user_rule",
				Weight:      0.90,
				Confidence:  1.0,
				Description: "Contains user-defined blacklisted term: '" + kwTrimmed + "'",
				MatchedText: kwTrimmed,
			})
		}
	}

	return signals
}
