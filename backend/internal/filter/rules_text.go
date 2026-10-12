package filter

import (
	"regexp"
	"strings"

	"aifreeyoutube/internal/models"
)

// LexicalRule associates a specific keyword/phrase pattern with a weight and diagnostic explanation.
type LexicalRule struct {
	Phrase      string
	Weight      float64
	Description string
}

// LLMPhrases lists distinctive words and clichés statistically overrepresented in LLM-generated output
// (e.g., standard ChatGPT prose defaults).
var LLMPhrases = []LexicalRule{
	{"delve into", 0.35, "Characteristic LLM prose marker ('delve into')"},
	{"delving into", 0.35, "Characteristic LLM prose marker ('delving into')"},
	{"in today's fast-paced world", 0.40, "Formulaic ChatGPT opening cliché"},
	{"in today's digital landscape", 0.40, "Formulaic ChatGPT opening cliché"},
	{"in today's digital world", 0.35, "Formulaic ChatGPT opening cliché"},
	{"tapestry of", 0.45, "Distinctive LLM metaphorical cliché ('tapestry of')"},
	{"testament to", 0.35, "Distinctive LLM essay phrasing ('testament to')"},
	{"beacon of", 0.35, "Distinctive LLM stock expression ('beacon of')"},
	{"game-changer", 0.25, "Overused synthetic marketing buzzword ('game-changer')"},
	{"game changer", 0.25, "Overused synthetic marketing buzzword ('game changer')"},
	{"revolutionize the way", 0.30, "Automated hype phrase ('revolutionize the way')"},
	{"unleash the power", 0.30, "Synthetic promotion phrasing ('unleash the power')"},
	{"look no further", 0.30, "Stock marketing script filler ('look no further')"},
	{"without further ado", 0.25, "Scripted filler introduction ('without further ado')"},
	{"embark on a journey", 0.35, "Formulaic robotic narrative opening ('embark on a journey')"},
	{"fascinating journey into", 0.35, "Formulaic robotic narration ('fascinating journey into')"},
	{"navigating the complexities", 0.35, "Standard LLM essay trope ('navigating the complexities')"},
	{"vital role in shaping", 0.30, "Academic-sounding automated filler ('vital role in shaping')"},
	{"take your productivity to the next level", 0.35, "Automated tech influencer script template"},
	{"in this comprehensive guide", 0.30, "Formulaic tutorial script header"},
	{"it's crucial to remember", 0.35, "ChatGPT ethical hedge phrasing ('crucial to remember')"},
	{"it is important to note", 0.30, "ChatGPT hedge formula ('important to note')"},
	{"step-by-step blueprint", 0.25, "Automated course/affiliate farm template"},
}

// ClickbaitFormulaRule associates a compiled regular expression with a weight and diagnostic explanation.
type ClickbaitFormulaRule struct {
	Regex       *regexp.Regexp
	Weight      float64
	Name        string
	Description string
}

// FormulaRules defines regex templates characteristic of automated faceless video factories.
var FormulaRules = []ClickbaitFormulaRule{
	{
		Regex:       regexp.MustCompile(`(?i)\bwhy nobody is talking about\b`),
		Weight:      0.40,
		Name:        "Formula: Why Nobody Is Talking About",
		Description: "Formulaic conspiracy/fear clickbait template engineered for automated channels",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bscientists just (discovered|found|made a terrifying|announced)\b`),
		Weight:      0.45,
		Name:        "Formula: Fake Science Alarmism",
		Description: "Synthetic science alarmism template common in automated stock-video mills",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\b\d+\s+(shocking|insane|terrifying|unbelievable|mind-blowing)\s+(things|facts|secrets|discoveries)\b`),
		Weight:      0.40,
		Name:        "Formula: Numbered Shock List",
		Description: "Automated listicle template for faceless channels",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bjust (changed|ruined|destroyed) everything forever\b`),
		Weight:      0.45,
		Name:        "Formula: Hyperbolic Catastrophism",
		Description: "Extreme hyperbolic clickbait typical of automated AI hype farms",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bthe disturbing (reality|truth|secret) (of|behind|about)\b`),
		Weight:      0.35,
		Name:        "Formula: Disturbing Reality Trope",
		Description: "Robotic drama/conspiracy narrative template",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bmake \$\d+[\d,kKmM]*\s*(per|a|\/)\s*(day|month|week)\b`),
		Weight:      0.60,
		Name:        "Formula: Automated Earning Claim",
		Description: "Formulaic earning promise commonly seen in automated scam mills",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\b(make|making|earn|earning) money (online )?(with|using) (ai|chatgpt|claude)\b`),
		Weight:      0.55,
		Name:        "Formula: AI Money Mill",
		Description: "Automated content formula generating AI money-making slop",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\b(laziest|easiest|simplest) way to make money (with|using|online)?\b`),
		Weight:      0.55,
		Name:        "Formula: Lazy AI Income Template",
		Description: "Formulaic hook for automated passive income channels",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\b(passive income|side hustle) (ideas )?(with|using) (ai|chatgpt|claude)\b`),
		Weight:      0.50,
		Name:        "Formula: AI Passive Income Scheme",
		Description: "Automated listicle template for passive income with AI bots",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bfaceless (ai )?(business|channel|youtube|video|influencer|model)\b`),
		Weight:      0.55,
		Name:        "Formula: Faceless Automation Farm",
		Description: "Content generated by or promoting automated faceless channel factories",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\b(paste|copy) this into (chatgpt|claude|ai)\b`),
		Weight:      0.60,
		Name:        "Formula: AI Prompt Farm Hook",
		Description: "Automated prompt-pasting clickbait template",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bai (model|influencer) (side hustle|that'?s making)\b`),
		Weight:      0.55,
		Name:        "Formula: AI Influencer Farm",
		Description: "Synthetic personality / virtual influencer farm template",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\bautomated (trading )?bot with ai\b`),
		Weight:      0.50,
		Name:        "Formula: Automated Bot Farm",
		Description: "Automated trading/crypto bot hype template",
	},
	{
		Regex:       regexp.MustCompile(`(?i)\b(chatgpt|ai) wrote this (script|video)\b`),
		Weight:      0.70,
		Name:        "Formula: Explicit AI Admission",
		Description: "Explicit declaration that script was authored by generative AI",
	},
}

// LeakedPromptMarkers checks for accidental system prompt leakage into titles or descriptions.
var LeakedPromptMarkers = []string{
	"as an ai language model",
	"certainly! here is a script",
	"certainly! here is a title",
	"sure, here's a compelling description",
	"act as a youtube creator",
	"generate a 10-minute script",
	"[insert channel name]",
	"[insert title here]",
}

// AITags lists tags and hashtags that explicitly identify synthetic media generation tools.
var AITags = map[string]float64{
	"#midjourney":       0.45,
	"#elevenlabs":       0.50,
	"#aivoice":          0.50,
	"#aiart":            0.35,
	"#aigenerated":      0.55,
	"#texttospeech":     0.45,
	"#tts":              0.35,
	"#sora":             0.40,
	"#runwayml":         0.40,
	"#stablediffusion":  0.40,
	"#chatgpt":          0.30,
	"#dalle":            0.35,
	"#synthesia":        0.50,
	"#heygen":           0.50,
	"#facelessyoutube":  0.50,
}

// EvaluateText inspects title and description for linguistic, structural, and metadata AI signals.
func EvaluateText(title, description string) []models.FilterSignal {
	var signals []models.FilterSignal

	combined := strings.ToLower(title + " " + description)
	titleLower := strings.ToLower(title)

	// 1. Check for blatant leaked prompts (Immediate high confidence trigger)
	for _, marker := range LeakedPromptMarkers {
		if strings.Contains(combined, marker) {
			signals = append(signals, models.FilterSignal{
				RuleName:    "Leaked LLM Prompt",
				Category:    "lexical",
				Weight:      0.85,
				Confidence:  0.95,
				Description: "Contains leaked system prompt or AI refusal text: '" + marker + "'",
				MatchedText: marker,
			})
		}
	}

	// 2. Check for LLM prose markers in title and description
	for _, rule := range LLMPhrases {
		if strings.Contains(titleLower, rule.Phrase) {
			// Title triggers receive a 1.25x boost because titles are concise
			signals = append(signals, models.FilterSignal{
				RuleName:    "LLM Phrasing (Title)",
				Category:    "lexical",
				Weight:      rule.Weight * 1.25,
				Confidence:  0.85,
				Description: rule.Description + " in video title",
				MatchedText: rule.Phrase,
			})
		} else if strings.Contains(combined, rule.Phrase) {
			signals = append(signals, models.FilterSignal{
				RuleName:    "LLM Phrasing (Description)",
				Category:    "lexical",
				Weight:      rule.Weight,
				Confidence:  0.75,
				Description: rule.Description,
				MatchedText: rule.Phrase,
			})
		}
	}

	// 3. Check for formulaic clickbait regex patterns
	for _, rule := range FormulaRules {
		if match := rule.Regex.FindString(title); match != "" {
			signals = append(signals, models.FilterSignal{
				RuleName:    rule.Name,
				Category:    "formula",
				Weight:      rule.Weight,
				Confidence:  0.90,
				Description: rule.Description,
				MatchedText: match,
			})
		}
	}

	// 4. Check for explicit AI generation tags and hashtags
	for tag, weight := range AITags {
		if strings.Contains(combined, tag) {
			signals = append(signals, models.FilterSignal{
				RuleName:    "Synthetic Media Tag",
				Category:    "metadata",
				Weight:      weight,
				Confidence:  0.95,
				Description: "Tagged with synthetic media generator: " + tag,
				MatchedText: tag,
			})
		}
	}

	return signals
}
