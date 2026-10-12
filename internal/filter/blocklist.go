package filter

import (
	"strings"
)

// KnownAIFarmPatterns contains canonical substrings or patterns frequently used
// by automated, faceless video factory channels.
var KnownAIFarmPatterns = []string{
	"daily facts ai",
	"mind blown ai",
	"scary stories tts",
	"reddit tales bot",
	"reddit stories tts",
	"history ai shorts",
	"luxury lifestyle ai",
	"motivational ai clips",
	"stoic ai wisdom",
	"ai recap",
	"ai movie recap",
	"automated facts",
	"faceless channel",
	"faceless youtube",
	"faceless video",
	"ai animated stories",
	"satisfying clips ai",
	"ai facts",
	"ai explained",
	"ai chronicles",
	"ai stories",
	"cash cow mastery",
	"chatgpt mastery",
	"reddit chill tales",
	"movie recap hd",
	"film recap official",
	"anime recap hd",
	"top 10 matrix",
	"mystery vault ai",
	"dark psychology hub",
	"sigma mindset daily",
}

// KnownAIFarmHandles contains confirmed mass-production bot handles.
var KnownAIFarmHandles = map[string]string{
	"deeplearningdaily_ai": "Automated AI news aggregator bot",
	"mindblowingfacts_ai":  "Synthetic voiceover stock clip farm",
	"redditstories_tts":    "Text-to-speech reddit bot farm",
	"historyuncovered_ai":  "AI hallucinated history documentary mill",
	"cryptonewstoday_ai":   "Automated crypto pump bot",
	"stoicismdaily_ai":     "ElevenLabs Marcus Aurelius quote generator",
	"movie_recaps_ai":      "Robotic narration movie synopsis mill",
	"luxurydiary_ai":       "Stock video luxury wealth automation",
}

// BlocklistManager handles fast querying for blocked and whitelisted entities.
type BlocklistManager struct{}

// NewBlocklistManager creates a blocklist manager.
func NewBlocklistManager() *BlocklistManager {
	return &BlocklistManager{}
}

// MatchKnownFarm checks if a channel display name or handle matches known automated farm patterns.
func (bm *BlocklistManager) MatchKnownFarm(channelName, channelHandle string) (bool, string) {
	nameLower := strings.ToLower(channelName)
	handleLower := strings.ToLower(channelHandle)

	// Direct match against known bot handles
	if reason, exists := KnownAIFarmHandles[handleLower]; exists {
		return true, reason
	}

	// Substring match against typical bot factory naming schemes
	for _, pattern := range KnownAIFarmPatterns {
		if strings.Contains(nameLower, pattern) || strings.Contains(handleLower, pattern) {
			return true, "Matches automated faceless channel naming convention: '" + pattern + "'"
		}
	}

	return false, ""
}

// IsCustomBlocked checks user-supplied channel IDs or handles.
func (bm *BlocklistManager) IsCustomBlocked(channelName, channelID, channelHandle string, customList []string) bool {
	nameLower := strings.ToLower(channelName)
	handleLower := strings.ToLower(channelHandle)
	idLower := strings.ToLower(channelID)

	for _, item := range customList {
		itemLower := strings.ToLower(strings.TrimSpace(item))
		if itemLower == "" {
			continue
		}
		if itemLower == idLower || itemLower == handleLower || strings.Contains(nameLower, itemLower) {
			return true
		}
	}
	return false
}

// IsWhitelisted checks if the creator has been explicitly marked as authentic human by the user.
func (bm *BlocklistManager) IsWhitelisted(channelID, channelHandle string, whitelist []string) bool {
	idLower := strings.ToLower(channelID)
	handleLower := strings.ToLower(channelHandle)

	for _, item := range whitelist {
		itemLower := strings.ToLower(strings.TrimSpace(item))
		if itemLower == "" {
			continue
		}
		if itemLower == idLower || itemLower == handleLower {
			return true
		}
	}
	return false
}
