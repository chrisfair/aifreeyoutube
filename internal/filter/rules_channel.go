package filter

import (
	"regexp"
	"strings"

	"aifreeyoutube/internal/models"
)

var (
	// SuspiciousHandleRegex matches bot handle patterns like "@HistoryFactsAI", "@StoicWisdom_ai", etc.
	SuspiciousHandleRegex = regexp.MustCompile(`(?i)(_ai$|^ai_|_official_ai|_gpt$|fact[sz]_ai|quotes_ai|recaps?_ai|automation_ai|tts_stories)`)

	// NumericSuffixHandleRegex matches channels created with automated naming schemes ending in 4+ digits
	NumericSuffixHandleRegex = regexp.MustCompile(`(?i)@[a-z]{3,}\d{4,}$`)

	// GeneratedBotHandleRegex matches standard YouTube auto-generated bot handles like "@user-ab3df7kz"
	GeneratedBotHandleRegex = regexp.MustCompile(`(?i)^user-[a-z0-9]{8,12}$`)

	// FactorySubjectTokens lists common subjects targeted by mass-production automated mills
	FactorySubjectTokens = []string{
		"fact", "facts", "scary", "horror", "creepy", "mystery", "mysteries",
		"crypto", "wealth", "money", "mindset", "stoic", "wisdom", "movie",
		"film", "anime", "recap", "recaps", "reddit", "motivation", "curiosity",
		"psychology", "dark", "sci", "science", "history", "space", "universe",
	}

	// FactoryNounTokens lists nouns commonly suffixed to create generic content mill names
	FactoryNounTokens = []string{
		"facts", "info", "compilation", "compilations", "recap", "recaps",
		"daily", "hub", "zone", "central", "vault", "archive", "tales",
		"lore", "matrix", "explained", "digest", "tidbits", "clips", "hq",
		"uncovered", "insider", "station", "corner", "chronicle", "chronicles",
	}
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
			Weight:      0.85,
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
			Weight:      0.65,
			Confidence:  0.90,
			Description: "Channel handle denotes automated AI generation: " + channelHandle,
			MatchedText: channelHandle,
		})
	} else if NumericSuffixHandleRegex.MatchString(channelHandle) {
		signals = append(signals, models.FilterSignal{
			RuleName:    "Algorithmic Username Pattern",
			Category:    "channel",
			Weight:      0.40,
			Confidence:  0.75,
			Description: "Channel handle utilizes automated keyword + numeric suffix format: " + channelHandle,
			MatchedText: channelHandle,
		})
	} else if GeneratedBotHandleRegex.MatchString(channelHandle) {
		signals = append(signals, models.FilterSignal{
			RuleName:    "Auto-Generated Bot Handle",
			Category:    "channel",
			Weight:      0.45,
			Confidence:  0.80,
			Description: "Channel utilizes unconfigured auto-generated bot handle identifier: " + channelHandle,
			MatchedText: channelHandle,
		})
	}

	// 4. Content Mill Channel Name Heuristic (e.g. "Sci Facts Info", "Chat History Compilations", "Daily Facts Hub")
	evaluateChannelNameHeuristics(channelName, &signals)

	return signals, false
}

// evaluateChannelNameHeuristics detects common factory naming patterns (Subject + Factory Noun).
func evaluateChannelNameHeuristics(channelName string, signals *[]models.FilterSignal) {
	nameLower := strings.ToLower(channelName)
	words := strings.Fields(nameLower)
	if len(words) < 2 {
		return
	}

	var hasSubject bool
	var hasNoun bool
	var matchedSubject string
	var matchedNoun string

	for _, w := range words {
		// Clean punctuation
		cleaned := strings.Trim(w, ",.!-#_")

		for _, s := range FactorySubjectTokens {
			if cleaned == s {
				hasSubject = true
				matchedSubject = s
				break
			}
		}

		for _, n := range FactoryNounTokens {
			if cleaned == n {
				hasNoun = true
				matchedNoun = n
				break
			}
		}
	}

	// If channel combines a factory subject with a factory noun (e.g. "Sci Facts Info", "History Compilations")
	if hasSubject && hasNoun && matchedSubject != matchedNoun {
		*signals = append(*signals, models.FilterSignal{
			RuleName:    "Content Mill Channel Formula",
			Category:    "channel",
			Weight:      0.55,
			Confidence:  0.85,
			Description: "Channel name matches automated mass-production factory naming convention (" + matchedSubject + " + " + matchedNoun + ")",
			MatchedText: channelName,
		})
	}
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
				Weight:      0.95,
				Confidence:  1.0,
				Description: "Contains user-defined blacklisted term: '" + kwTrimmed + "'",
				MatchedText: kwTrimmed,
			})
		}
	}

	return signals
}
