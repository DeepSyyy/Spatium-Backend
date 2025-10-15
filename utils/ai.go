package utils

import (
	"context"
	"fmt"

	"github.com/DeepSyyy/Spatium-Backend/config"
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
