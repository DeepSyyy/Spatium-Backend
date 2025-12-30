package utils

import (
	"context"
	"regexp"
	"strings"

	"github.com/DeepSyyy/Spatium-Backend/config"
	openai "github.com/sashabaranov/go-openai"
)

// ModerationResult contains the result of content moderation
type ModerationResult struct {
	IsSafe          bool     `json:"is_safe"`
	IsBlocked       bool     `json:"is_blocked"`
	IsCrisisContent bool     `json:"is_crisis_content"` // New: flag for crisis content
	Categories      []string `json:"categories,omitempty"`
	Reason          string   `json:"reason,omitempty"`
	SeverityLevel   string   `json:"severity_level,omitempty"`  // "low", "medium", "high", "critical"
	SupportMessage  string   `json:"support_message,omitempty"` // Crisis support message
}

// TriggerWords - kata-kata yang memerlukan perhatian khusus untuk kesehatan mental
// PENTING: Konten dengan kata-kata ini akan DIBLOKIR dari publik untuk mencegah Werther effect
var triggerWords = []string{
	// Suicide-related (Indonesian) - termasuk variasi slang/typo
	"bunuh diri", "bunuhdiri", "bundir", "bnuh diri", "bunuh dir",
	"mau mati", "maumati", "mw mati", "mo mati", "pgn mati",
	"ingin mati", "pengen mati", "pngen mati", "pengin mati",
	"lebih baik mati", "mending mati", "mendingan mati",
	"gantung diri", "gantungdiri", "gtung diri",
	"akhiri hidup", "akhirin hidup", "end it all",
	"tidak ingin hidup", "ga mau hidup", "gak mau hidup", "gamau hidup",
	"menyerah hidup", "nyerah hidup", "give up",
	"capek hidup", "cape hidup", "cpk hidup", "lelah hidup",
	"bosan hidup", "bosen hidup", "muak hidup",
	"mati aja", "mati saja", "matiin aja", "matiaja",
	"ga ada gunanya hidup", "hidup ga ada artinya",
	"loncat dari", "lompat dari", "terjun dari",
	// Self-harm (Indonesian) - termasuk variasi
	"menyakiti diri", "nyakitin diri", "sakitin diri sendiri",
	"lukai diri", "melukai diri", "luka diri",
	"potong tangan", "potong nadi", "iris tangan", "iris nadi",
	"silet", "cutter", "sayat", "nyayat",
	"self harm", "selfharm", "sh",
	// Suicide-related (English)
	"kill myself", "kms", "kys", "kill yourself",
	"want to die", "wanna die", "i wanna die",
	"end my life", "end it all", "ending it",
	"suicide", "suicidal", "commit suicide",
	"hang myself", "hanging myself",
	"jump off", "jumping off",
	"overdose", "od",
	// Self-harm (English)
	"cut myself", "cutting myself", "cutting",
	"hurt myself", "hurting myself",
}

// Crisis patterns - regex patterns untuk deteksi lebih fleksibel
var crisisPatterns = []string{
	`(?i)pen?gen\s*(bunuh|mati|bundir)`,
	`(?i)(mau|ingin|pengen)\s*(mati|bundir|bunuh\s*diri)`,
	`(?i)(cape[k]?|lelah|bosan|bosen|muak)\s*(hidup|sama\s*hidup)`,
	`(?i)(akhiri|selesaikan|end)\s*(hidup|semuanya|it\s*all)`,
	`(?i)(loncat|lompat|terjun)\s*(dari|ke)\s*(gedung|jembatan|lantai)`,
	`(?i)(iris|potong|sayat)\s*(tangan|nadi|pergelangan)`,
	`(?i)(ga[k]?|tidak|no)\s*(ada|punya)\s*(harapan|hope|gunanya)`,
	`(?i)(lebih\s*baik|mending)\s*(ga[k]?\s*ada|mati|pergi)`,
	`(?i)i\s*(want|wanna)\s*(to\s*)?(die|kill\s*myself)`,
}

// ToxicWords - kata-kata kasar/toxic yang harus diblokir
var toxicWords = []string{
	// Indonesian profanity/toxic
	"anjing", "bangsat", "babi", "goblok", "tolol", "idiot",
	"bego", "bodoh", "kampret", "bajingan", "tai", "kontol",
	"memek", "ngentot", "asu", "jancok", "cuk", "bacot",
	// Bullying words (Indonesian)
	"jelek banget", "mati aja", "sampah", "ga guna", "tidak berguna",
	"pembawa sial", "kutukan", "menyedihkan",
	// English profanity
	"fuck", "shit", "bitch", "asshole", "bastard", "dick",
	"cunt", "whore", "slut", "retard", "faggot",
	// Bullying words (English)
	"loser", "worthless", "pathetic", "disgusting", "kill yourself",
}

// ModerateContent checks content for harmful or toxic material
// Returns ModerationResult with safety status and categories
func ModerateContent(content string) (*ModerationResult, error) {
	result := &ModerationResult{
		IsSafe:          true,
		IsBlocked:       false,
		IsCrisisContent: false,
		Categories:      []string{},
		SeverityLevel:   "low",
	}

	lowerContent := strings.ToLower(content)

	// Step 1: Check for crisis indicators (BLOCK from public + show support)
	// Untuk mencegah Werther effect, konten krisis tidak dipublikasikan
	for _, word := range triggerWords {
		if strings.Contains(lowerContent, word) {
			result.IsSafe = false
			result.IsBlocked = true
			result.IsCrisisContent = true
			result.Categories = append(result.Categories, "crisis_indicator")
			result.SeverityLevel = "critical"
			result.Reason = "Konten mengandung indikator krisis kesehatan mental"
			result.SupportMessage = GetCrisisSupportMessage()
			return result, nil
		}
	}

	// Step 1b: Check crisis patterns with regex
	for _, pattern := range crisisPatterns {
		matched, _ := regexp.MatchString(pattern, lowerContent)
		if matched {
			result.IsSafe = false
			result.IsBlocked = true
			result.IsCrisisContent = true
			result.Categories = append(result.Categories, "crisis_indicator")
			result.SeverityLevel = "critical"
			result.Reason = "Konten mengandung indikator krisis kesehatan mental"
			result.SupportMessage = GetCrisisSupportMessage()
			return result, nil
		}
	}

	// Step 2: Check for toxic/profane words
	for _, word := range toxicWords {
		if strings.Contains(lowerContent, word) {
			result.IsSafe = false
			result.IsBlocked = true
			result.Categories = append(result.Categories, "toxic_language")
			result.Reason = "Konten mengandung bahasa yang tidak pantas"
			result.SeverityLevel = "high"
			return result, nil
		}
	}

	// Step 3: Check for hate speech patterns
	hatePatterns := []string{
		`(?i)(semua|seluruh)\s+(orang|manusia)\s+\w+\s+(harus|layak)\s+(mati|dibunuh)`,
		`(?i)(bunuh|hajar|siksa)\s+(semua|seluruh)`,
		`(?i)(ras|suku|agama)\s+\w+\s+(sampah|harus\s+mati)`,
	}

	for _, pattern := range hatePatterns {
		matched, _ := regexp.MatchString(pattern, lowerContent)
		if matched {
			result.IsSafe = false
			result.IsBlocked = true
			result.Categories = append(result.Categories, "hate_speech")
			result.Reason = "Konten mengandung ujaran kebencian"
			result.SeverityLevel = "critical"
			return result, nil
		}
	}

	// Step 4: Use OpenAI Moderation API for additional check (if available)
	if config.AppConfig != nil && config.AppConfig.OpenAIAPIKey != "" {
		openAIResult, err := checkWithOpenAIModeration(content)
		if err == nil && openAIResult != nil {
			if openAIResult.IsBlocked {
				result.IsSafe = false
				result.IsBlocked = true
				result.Categories = append(result.Categories, openAIResult.Categories...)
				result.Reason = openAIResult.Reason
				if openAIResult.SeverityLevel == "critical" || openAIResult.SeverityLevel == "high" {
					result.SeverityLevel = openAIResult.SeverityLevel
				}
			}
		}
	}

	return result, nil
}

// checkWithOpenAIModeration uses OpenAI's moderation endpoint
func checkWithOpenAIModeration(content string) (*ModerationResult, error) {
	client := openai.NewClient(config.AppConfig.OpenAIAPIKey)
	ctx := context.Background()

	resp, err := client.Moderations(ctx, openai.ModerationRequest{
		Input: content,
	})
	if err != nil {
		return nil, err
	}

	result := &ModerationResult{
		IsSafe:        true,
		IsBlocked:     false,
		Categories:    []string{},
		SeverityLevel: "low",
	}

	if len(resp.Results) > 0 {
		modResult := resp.Results[0]

		if modResult.Flagged {
			result.IsSafe = false
			result.IsBlocked = true
			result.SeverityLevel = "high"

			// Map OpenAI categories
			if modResult.Categories.Hate {
				result.Categories = append(result.Categories, "hate")
				result.SeverityLevel = "critical"
			}
			if modResult.Categories.HateThreatening {
				result.Categories = append(result.Categories, "hate_threatening")
				result.SeverityLevel = "critical"
			}
			if modResult.Categories.Harassment {
				result.Categories = append(result.Categories, "harassment")
			}
			if modResult.Categories.HarassmentThreatening {
				result.Categories = append(result.Categories, "harassment_threatening")
				result.SeverityLevel = "critical"
			}
			if modResult.Categories.SelfHarm {
				result.Categories = append(result.Categories, "self_harm")
				result.SeverityLevel = "critical"
			}
			if modResult.Categories.SelfHarmIntent {
				result.Categories = append(result.Categories, "self_harm_intent")
				result.SeverityLevel = "critical"
			}
			if modResult.Categories.SelfHarmInstructions {
				result.Categories = append(result.Categories, "self_harm_instructions")
				result.SeverityLevel = "critical"
			}
			if modResult.Categories.Sexual {
				result.Categories = append(result.Categories, "sexual")
			}
			if modResult.Categories.Violence {
				result.Categories = append(result.Categories, "violence")
			}
			if modResult.Categories.ViolenceGraphic {
				result.Categories = append(result.Categories, "violence_graphic")
				result.SeverityLevel = "critical"
			}

			result.Reason = "Konten melanggar pedoman komunitas"
		}
	}

	return result, nil
}

// SanitizeContent removes or replaces harmful content
// This is a lighter approach - replaces bad words with asterisks
func SanitizeContent(content string) string {
	result := content

	for _, word := range toxicWords {
		// Create asterisk replacement of same length
		replacement := strings.Repeat("*", len(word))
		// Case-insensitive replacement
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(word))
		result = re.ReplaceAllString(result, replacement)
	}

	return result
}

// ContainsCrisisIndicators checks if content indicates mental health crisis
// This is for routing to appropriate support resources, not blocking
func ContainsCrisisIndicators(content string) bool {
	lowerContent := strings.ToLower(content)

	for _, word := range triggerWords {
		if strings.Contains(lowerContent, word) {
			return true
		}
	}

	return false
}

// GetCrisisSupportMessage returns a supportive message for crisis situations
func GetCrisisSupportMessage() string {
	return `Kami mendeteksi bahwa kamu mungkin sedang mengalami masa sulit. 
Ingat, kamu tidak sendirian. Berikut adalah sumber bantuan yang bisa kamu hubungi:

📞 Into The Light Indonesia: 119 ext 8
📞 Yayasan Pulih: (021) 788-42580
📞 LSM Jangan Bunuh Diri: 021-9696 9293 / 0858-8000-0123

Kami peduli dengan kesejahteraanmu. 💙`
}
