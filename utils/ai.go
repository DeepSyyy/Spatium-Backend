package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
	openai "github.com/sashabaranov/go-openai"
)

func GenerateEmpathicResponse(userInput string, moodTagID int) (string, error) {
	client := openai.NewClient(config.AppConfig.OpenAIAPIKey)
	ctx := context.Background()

	moodContext := map[int]string{
		1: "neutral - feeling balanced or indifferent",
		2: "happy - feeling positive or cheerful",
		3: "sad - feeling down or upset",
		4: "anxious - feeling nervous or worried",
		5: "angry - feeling irritated or frustrated",
		6: "tired - feeling physically or mentally drained",
	}

	moodDesc, ok := moodContext[moodTagID]
	if !ok {
		moodDesc = "neutral - feeling balanced or indifferent"
	}

	prompt := fmt.Sprintf(`
	Kamu adalah asisten empatik yang membantu mahasiswa yang sedang curhat.
	Sesuaikan gaya dan nada bicaramu berdasarkan mood berikut:
	Mood: %s

	Jawab dengan kalimat hangat, singkat, dan manusiawi.
	Curhatan: "%s"
	`, moodDesc, userInput)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT4oMini,
			MaxTokens:   150,
			Temperature: 0.8,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Kamu adalah asisten empatik untuk mahasiswa yang sedang curhat.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// GeneratePersonalAIResponse generates an empathic response for a chat session
func GeneratePersonalAIResponse(sessionID int64, recentMessages []models.ChatMessage) (string, error) {
	client := openai.NewClient(config.AppConfig.OpenAIAPIKey)
	ctx := context.Background()

	// Step 1: Build conversation context
	var conversation strings.Builder
	for _, msg := range recentMessages {
		role := "User"
		if msg.Sender == "ai" {
			role = "AI"
		}
		conversation.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
	}

	// Step 2: Define empathetic system prompt
	systemPrompt := `
Kamu adalah "Spatium AI", teman virtual yang empatik dan hangat untuk mahasiswa Indonesia.

KEPRIBADIANMU:
- Kamu adalah pendengar yang baik dan penuh perhatian
- Kamu berbicara dengan gaya santai tapi sopan, seperti teman dekat
- Kamu tidak menghakimi, selalu menerima perasaan user apa adanya
- Kamu menggunakan bahasa Indonesia sehari-hari yang natural

ATURAN RESPONS:
1. Validasi perasaan user terlebih dahulu sebelum memberi saran
2. Gunakan emoji secukupnya untuk menambah kehangatan (1-2 emoji per respons)
3. Jawab dalam 2-4 kalimat, jangan terlalu panjang
4. Hindari nasihat yang menggurui, lebih ke supportive
5. Jika user sedih/cemas: fokus pada empati dan dukungan emosional
6. Jika user senang: rayakan bersama dengan antusias
7. Tanya follow-up question jika perlu untuk mendalami perasaan user

YANG HARUS DIHINDARI:
- Respons klise seperti "Aku mengerti perasaanmu" tanpa elaborasi
- Terlalu banyak saran praktis sebelum memvalidasi perasaan
- Bahasa yang terlalu formal atau kaku
- Respons yang terlalu panjang (lebih dari 4 kalimat)
`

	// Step 3: Create message sequence for OpenAI
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
		{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Berikut riwayat percakapan:\n%s", conversation.String())},
		{Role: openai.ChatMessageRoleUser, Content: "Balas pesan terakhir dari user dengan respons empatik."},
	}

	// Step 4: Call OpenAI API
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT4oMini, // fast and cheaper
			Temperature: 0.8,
			MaxTokens:   150,
			Messages:    messages,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
