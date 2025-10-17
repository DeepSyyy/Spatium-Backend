package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatSession struct {
	InternalID int64         `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID     `json:"public_id" db:"public_id"`
	UserID     int64         `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	Title      string        `json:"title" db:"title" gorm:"type:varchar(255);not null"`
	MoodTag    string        `json:"mood_tag" db:"mood_tag" gorm:"type:varchar(50)"`
	CreatedAt  time.Time     `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time     `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
	Messages   []ChatMessage `json:"messages" gorm:"foreignKey:SessionID;references:InternalID"`
}

type ChatSessionResponse struct {
	PublicID  string `json:"public_id"`
	Title     string `json:"title"`
	MoodTag   string `json:"mood_tag"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
