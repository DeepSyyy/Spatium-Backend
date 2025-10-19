package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type ChatSessionRepository interface {
	Create(session *models.ChatSession) error
	FindByUserID(userID int64) ([]models.ChatSession, error)
	FindByPublicID(publicID string) (*models.ChatSession, error)
	FindByPublicIDAndUserID(publicID string, userID int64) (*models.ChatSession, error)
	Delete(sessionID string) error
}

type chatSessionRepository struct {
	db *gorm.DB
}

func NewChatSessionRepository() ChatSessionRepository {
	return &chatSessionRepository{config.DB}
}

func (r *chatSessionRepository) Create(session *models.ChatSession) error {
	return r.db.Create(session).Error
}

func (r *chatSessionRepository) FindByUserID(userID int64) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := r.db.
		Where("user_internal_id = ?", userID).
		Order("created_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatSessionRepository) FindByPublicID(publicID string) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.
		Where("public_id = ?", publicID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *chatSessionRepository) Delete(sessionID string) error {
	return r.db.
		Where("public_id = ?", sessionID).
		Delete(&models.ChatSession{}).Error
}

func (r *chatSessionRepository) FindByPublicIDAndUserID(publicID string, userID int64) (*models.ChatSession, error) {
	var session models.ChatSession
	err := r.db.
		Where("public_id = ? AND user_internal_id = ?", publicID, userID).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
