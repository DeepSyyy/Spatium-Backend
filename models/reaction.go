package models

import (
	"time"

	"github.com/google/uuid"
)

type Reaction struct {
	InternalID             int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID               uuid.UUID `json:"public_id" db:"public_id"`
	UserID                 int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	PostID                 int64     `json:"post_internal_id" db:"post_internal_id" gorm:"column:post_internal_id;not null"`
	ReactionTypeInternalID int64     `json:"reaction_type_internal_id" gorm:"not null"`
	Emoji                  string    `json:"emoji" db:"emoji" gorm:"type:varchar(10);not null"`
	CreatedAt              time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
}

type ReactionResponse struct {
	Emoji string `json:"emoji"`
	Count int64  `json:"count"`
}
