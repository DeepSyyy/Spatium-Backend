package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	InternalID  int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID    uuid.UUID `json:"public_id" db:"public_id"`
	UserID      int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	UserMessage string    `json:"user_message" db:"user_message" gorm:"type:text;not null"`
	AiMessage   string    `json:"ai_message" db:"ai_message" gorm:"type:text;not null"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
}
