package services

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
)

type ChatSessionService interface {
	Create(userID int64, title string) (*models.ChatSession, error)
	GetUserSessions(userID int64) ([]models.ChatSession, error)
	GetSessionByPublicID(publicID string) (*models.ChatSession, error)
	DeleteSession(sessionID string) error
}

type chatSessionService struct {
	repo repositories.ChatSessionRepository
}

func NewChatSessionService(repo repositories.ChatSessionRepository) ChatSessionService {
	return &chatSessionService{repo}
}

func (s *chatSessionService) Create(userID int64, title string) (*models.ChatSession, error) {
	session := &models.ChatSession{
		PublicID:  uuid.New(),
		UserID:    userID,
		Title:     title,
		MoodTag:   "",
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *chatSessionService) GetUserSessions(userID int64) ([]models.ChatSession, error) {
	return s.repo.FindByUserID(userID)
}

func (s *chatSessionService) GetSessionByPublicID(publicID string) (*models.ChatSession, error) {
	return s.repo.FindByPublicID(publicID)
}

func (s *chatSessionService) DeleteSession(sessionID string) error {
	return s.repo.Delete(sessionID)
}
