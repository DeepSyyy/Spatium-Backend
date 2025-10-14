package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	UserID     int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	AiResponse string    `json:"ai_response" db:"ai_response"`
	Content    string    `json:"content" db:"content" gorm:"type:text;not null"`
	MoodTagID  int64     `json:"mood_tag_internal_id" db:"mood_tag_internal_id" gorm:"column:mood_tag_internal_id;not null"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
}
