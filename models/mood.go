package models

import (
	"time"

	"github.com/google/uuid"
)

type Mood struct {
	InternalID   int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID     uuid.UUID `json:"public_id" db:"public_id"`
	Emoji        string    `json:"emoji" db:"emoji"`
	Note         string    `json:"note" db:"note"`
	UserID       int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	AiReflection string    `json:"ai_reflection" db:"ai_reflection"`
	Date         time.Time `json:"date" db:"date" gorm:"autoCreateTime"`
}
