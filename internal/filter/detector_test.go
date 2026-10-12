package filter

import (
	"testing"

	"aifreeyoutube/internal/models"
)

func TestDetectorScoring(t *testing.T) {
	configMgr := NewConfigManager()
	detector := NewDetector(configMgr)
	settings := configMgr.GetSettings()

	tests := []struct {
		name          string
		video         models.Video
		shouldPass    bool
		minAIScore    float64
		maxAIScore    float64
		expectVerdict models.FilterVerdict
	}{
		{
			name: "Authentic Human Content",
			video: models.Video{
				ID:          "human1",
				Title:       "Building a Custom Mechanical Keyboard From Scratch",
				ChannelName: "Handmade Crafts",
				Description: "Soldering switches, lubing stabilizers, and testing keycaps.",
				Duration:    "18:42",
				IsShort:     false,
			},
			shouldPass:    true,
			minAIScore:    0.0,
			maxAIScore:    0.20,
			expectVerdict: models.VerdictClean,
		},
		{
			name: "Blatant ChatGPT Phrasing & Slop",
			video: models.Video{
				ID:          "ai1",
				Title:       "Delve Into The Tapestry Of Future Artificial Intelligence",
				ChannelName: "Tech News Daily",
				Description: "In today's fast-paced world, this game-changer will revolutionize everything. Without further ado, let us embark on a journey.",
				Duration:    "10:05",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.60,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Clickbait Alarmist Science Slop",
			video: models.Video{
				ID:          "ai2",
				Title:       "Scientists Just Discovered Terrifying Anomaly In Deep Space",
				ChannelName: "Daily Facts AI",
				Description: "Why nobody is talking about this shocking discovery that changes everything forever.",
				Duration:    "08:14",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.65,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Leaked Prompt Slop",
			video: models.Video{
				ID:          "ai3",
				Title:       "Best Productivity Hacks For Students",
				ChannelName: "Student Guide",
				Description: "Certainly! Here is a script for a youtube video highlighting top study tips.",
				Duration:    "07:30",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.70,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Listicle Fact Mill With Channel Formula",
			video: models.Video{
				ID:          "fact1",
				Title:       "100 History Facts | The Strangest History Facts You'll Ever Hear !!",
				ChannelName: "Sci Facts Info",
				Description: "Strange history facts compilation.",
				Duration:    "15:30",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.50,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Faceless Cash Cow Channel Automation",
			video: models.Video{
				ID:          "cashcow1",
				Title:       "Faceless YouTube Channel Automation: Make $10,000/Month with AI",
				ChannelName: "Passive Income Hub",
				Description: "How to make money with ChatGPT and Elevenlabs on autopilot.",
				Duration:    "11:20",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.75,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Reddit Scraping TTS Story Farm",
			video: models.Video{
				ID:          "reddit1",
				Title:       "r/AskReddit - What is the most terrifying thing you have ever experienced?",
				ChannelName: "Reddit Tales",
				Description: "Top Reddit stories read aloud with text to speech.",
				Duration:    "24:00",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.70,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Fake Science Space Alarmism",
			video: models.Video{
				ID:          "space1",
				Title:       "James Webb Telescope Just Detected Something Terrifying In Deep Space",
				ChannelName: "Universe Mysteries",
				Description: "Scientists are terrified by what was found inside.",
				Duration:    "08:45",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.65,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Movie Recap Farm",
			video: models.Video{
				ID:          "recap1",
				Title:       "Movie Recap: Man Gets Trapped In An Underground Bunker For 100 Years",
				ChannelName: "Movie Recap Zone",
				Description: "Plot explained in 12 minutes.",
				Duration:    "12:00",
				IsShort:     false,
			},
			shouldPass:    false,
			minAIScore:    0.65,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Shorts Elimination",
			video: models.Video{
				ID:          "short1",
				Title:       "Satisfying Soap Cutting #shorts",
				ChannelName: "Clip Master",
				Description: "Satisfying compilation #shorts #asmr",
				Duration:    "0:45",
				IsShort:     true,
			},
			shouldPass:    false,
			minAIScore:    0.80,
			maxAIScore:    1.0,
			expectVerdict: models.VerdictAISlop,
		},
		{
			name: "Whitelisted Creator Override",
			video: models.Video{
				ID:          "white1",
				Title:       "Why Nobody Is Talking About This Tech Paradigm Shift",
				ChannelName: "Trusted Human Tech",
				ChannelID:   "UC_trusted_human",
				Description: "A deep dive into computational complexity.",
				Duration:    "22:15",
				IsShort:     false,
			},
			shouldPass:    true,
			minAIScore:    0.0,
			maxAIScore:    0.0,
			expectVerdict: models.VerdictClean,
		},
	}

	// Add whitelist entry for test case
	settings.WhitelistedChannels = []string{"UC_trusted_human"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated, passed := detector.Evaluate(&tt.video, &settings)

			if passed != tt.shouldPass {
				t.Errorf("expected passed=%v, got %v (AIScore: %f, Verdict: %s, Signals: %d)",
					tt.shouldPass, passed, evaluated.AIScore, evaluated.Verdict, len(evaluated.Signals))
			}

			if evaluated.AIScore < tt.minAIScore || evaluated.AIScore > tt.maxAIScore {
				t.Errorf("expected AIScore between %f and %f, got %f",
					tt.minAIScore, tt.maxAIScore, evaluated.AIScore)
			}

			if evaluated.Verdict != tt.expectVerdict {
				t.Errorf("expected Verdict=%s, got %s", tt.expectVerdict, evaluated.Verdict)
			}
		})
	}
}
