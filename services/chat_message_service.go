package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/DeepSyyy/Spatium-Backend/utils"
)

type ChatMessageService interface {
	Create(messageContent string, sessionID int64, sender string) (*models.ChatMessage, error)
	FindBySession(sessionID int64, limit int) ([]models.ChatMessage, error)
	FindLastMessages(sessionID int64, limit int) ([]models.ChatMessage, error)
	ConvertPublicIDToInternalIDSession(publicID string) (int64, error)
	ValidateSessionOwnership(userID int64, sessionPublicID string) (int64, error)
}

type chatMessageService struct {
	repoMessage repositories.ChatMessageRepository
	repoSession repositories.ChatSessionRepository
}

func NewChatMessageService(repoM repositories.ChatMessageRepository, repoS repositories.ChatSessionRepository) ChatMessageService {
	return &chatMessageService{repoM, repoS}
}

func (s *chatMessageService) Create(messageContent string, sessionID int64, sender string) (*models.ChatMessage, error) {
	if messageContent == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if sender != "user" && sender != "ai" {
		sender = "user" // fallback safety
	}

	// 1️⃣ Simpan pesan user
	message := &models.ChatMessage{
		SessionID: sessionID,
		Content:   messageContent,
		Sender:    sender,
		CreatedAt: time.Now(),
	}

	if err := s.repoMessage.Create(message); err != nil {
		return nil, err
	}

	// 2️⃣ Jika pengirim adalah user, AI akan merespons otomatis
	if sender == "user" {
		// Ambil pesan terakhir (misal 6 untuk konteks)
		lastMessages, err := s.repoMessage.FindLastMessages(sessionID, 6)
		if err != nil {
			fmt.Println("⚠️ Failed to get context messages:", err)
			return message, nil // tetap return pesan user meski AI gagal
		}

		// 3️⃣ Generate response dari AI
		aiResponse, err := utils.GeneratePersonalAIResponse(sessionID, lastMessages)
		if err != nil {
			fmt.Println("⚠️ AI response failed:", err)
			return message, nil
		}

		// 4️⃣ Simpan pesan AI juga
		aiMessage := &models.ChatMessage{
			SessionID: sessionID,
			Content:   aiResponse,
			Sender:    "ai",
			CreatedAt: time.Now(),
		}
		_ = s.repoMessage.Create(aiMessage)
	}

	// 5️⃣ Return pesan user (AI message bisa diambil via GET endpoint)
	return message, nil
}

func (s *chatMessageService) FindBySession(sessionID int64, limit int) ([]models.ChatMessage, error) {
	return s.repoMessage.FindBySession(sessionID, limit)
}

func (s *chatMessageService) FindLastMessages(sessionID int64, limit int) ([]models.ChatMessage, error) {
	return s.repoMessage.FindLastMessages(sessionID, limit)
}

func (s *chatMessageService) ConvertPublicIDToInternalIDSession(publicID string) (int64, error) {
	return s.repoMessage.ConvertPublicIDToInternalIDSession(publicID)
}

func (s *chatMessageService) ValidateSessionOwnership(userID int64, sessionPublicID string) (int64, error) {
	session, err := s.repoSession.FindByPublicIDAndUserID(sessionPublicID, userID)
	if err != nil {
		return 0, fmt.Errorf("unauthorized: session not found or does not belong to you")
	}
	return session.InternalID, nil
}
