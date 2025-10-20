package models

import (
	"time"

	"github.com/google/uuid"
)

type AIReflection struct {
	InternalID int64     `json:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" gorm:"uniqueIndex;not null"`
	UserID     int64     `json:"user_internal_id" gorm:"column:user_internal_id;not null"`
	MoodTagID  int64     `json:"mood_tag_internal_id" gorm:"column:mood_tag_internal_id;not null"`
	Reflection string    `json:"reflection" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type AIReflectionResponse struct {
	PublicID   string    `json:"public_id"`
	MoodTagID  int64     `json:"mood_tag_internal_id"`
	Reflection string    `json:"reflection"`
	CreatedAt  time.Time `json:"created_at"`
}
