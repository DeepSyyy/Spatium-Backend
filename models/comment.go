package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	UserID     int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	PostID     int64     `json:"post_internal_id" db:"post_internal_id" gorm:"column:post_internal_id;not null"`
	Content    string    `json:"content" db:"content" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
}

type CommentResponse struct {
	PublicID  string    `json:"public_id"`
	UserAlias string    `json:"alias,omitempty"` // opsional: tampilkan alias user
	Content   string    `json:"content"`
	PostID    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}
