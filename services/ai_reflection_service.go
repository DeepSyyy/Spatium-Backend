package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type AIReflectionService interface {
	GenerateReflection(userID, moodTagID int64, moodDesc string) (*models.AIReflection, error)
	GetUserReflections(userID int64) ([]models.AIReflection, error)
}

type aiReflectionService struct {
	repo repositories.AIReflectionRepository
}

func NewAIReflectionService(repo repositories.AIReflectionRepository) AIReflectionService {
	return &aiReflectionService{repo}
}

func (s *aiReflectionService) GenerateReflection(userID, moodTagID int64, moodDesc string) (*models.AIReflection, error) {
	client := openai.NewClient(config.AppConfig.OpenAIAPIKey)
	ctx := context.Background()

	prompt := fmt.Sprintf(`
	Kamu adalah asisten empatik yang membantu mahasiswa merefleksikan perasaannya hari ini.
	Mood pengguna: %s.
	Berikan refleksi singkat dan hangat (1–3 kalimat) untuk membantu mereka memahami perasaannya.
	`, moodDesc)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: "Kamu adalah asisten refleksi empatik."},
				{Role: openai.ChatMessageRoleUser, Content: prompt},
			},
			MaxTokens:   100,
			Temperature: 0.7,
		},
	)
	if err != nil {
		return nil, err
	}

	reflection := &models.AIReflection{
		UserID:     userID,
		PublicID:   uuid.New(),
		MoodTagID:  moodTagID,
		Reflection: resp.Choices[0].Message.Content,
		CreatedAt:  time.Now(),
	}

	err = s.repo.Create(reflection)
	return reflection, err
}

func (s *aiReflectionService) GetUserReflections(userID int64) ([]models.AIReflection, error) {
	return s.repo.GetByUser(userID)
}
