package models

import (
	"time"

	"github.com/google/uuid"
)

type DailyMood struct {
	InternalID        int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID          uuid.UUID `json:"public_id" db:"public_id" gorm:"type:uuid;default:gen_random_uuid()"`
	UserInternalID    int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	MoodTagInternalID int64     `json:"mood_tag_internal_id" db:"mood_tag_internal_id" gorm:"column:mood_tag_internal_id;not null"`
	Note              string    `json:"note" db:"note" gorm:"type:text"`
	Date              time.Time `json:"date" db:"date" gorm:"type:date;not null"`
	CreatedAt         time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

type DailyMoodResponse struct {
	PublicID  string    `json:"public_id"`
	MoodTagID int64     `json:"mood_tag_internal_id"`
	Note      string    `json:"note"`
	Date      string    `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
