package models

import (
	"time"

	"github.com/google/uuid"
)

type ReactionType struct {
	InternalID  int64     `json:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID    uuid.UUID `json:"public_id" gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(50);unique;not null"`
	Emoji       string    `json:"emoji" gorm:"type:varchar(10);not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
