package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"gorm.io/gorm"
)

type ChatMessageRepository interface {
	Create(message *models.ChatMessage) error
	FindBySession(sessionID int64, limit int) ([]models.ChatMessage, error)
	FindLastMessages(sessionID int64, limit int) ([]models.ChatMessage, error)
	ConvertPublicIDToInternalIDSession(publicID string) (int64, error)
}

type chatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository() ChatMessageRepository {
	return &chatMessageRepository{config.DB}
}

func (r *chatMessageRepository) Create(message *models.ChatMessage) error {
	encrypted, err := utils.Encrypt(message.Content)
	if err != nil {
		return err
	}
	message.Content = encrypted
	return r.db.Create(message).Error
}

func (r *chatMessageRepository) FindBySession(sessionID int64, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.
		Where("session_internal_id = ?", sessionID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	res, err := decryptMessages(messages)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *chatMessageRepository) FindLastMessages(sessionID int64, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.
		Where("session_internal_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	messages, err = decryptMessages(messages)
	if err != nil {
		return nil, err
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func decryptMessages(messages []models.ChatMessage) ([]models.ChatMessage, error) {
	for i := range messages {
		decrypted, err := utils.Decrypt(messages[i].Content)
		if err != nil {
			return nil, err
		}
		messages[i].Content = decrypted
	}
	return messages, nil
}

func (r *chatMessageRepository) ConvertPublicIDToInternalIDSession(publicID string) (int64, error) {
	var session models.ChatSession
	err := r.db.Select("internal_id").Where("public_id = ?", publicID).First(&session).Error
	if err != nil {
		return 0, err
	}
	return session.InternalID, nil
}
