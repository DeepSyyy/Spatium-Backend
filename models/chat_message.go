package models

import "time"

type ChatMessage struct {
	InternalID int64  `json:"-" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	SessionID  int64  `json:"session_internal_id" db:"session_internal_id" gorm:"column:session_internal_id;not null"`
	Sender     string `json:"sender" db:"sender" gorm:"type:varchar(50);not null"` // "user" or "ai"
	Content    string `json:"content" db:"content" gorm:"type:text;not null"`
	// MoodTag    string       `json:"mood_tag" db:"mood_tag" gorm:"type:varchar(50)"`
	CreatedAt time.Time    `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	Session   *ChatSession `json:"-" gorm:"foreignKey:SessionID;references:InternalID"`
}

type ChatMessageResponse struct {
	SessionID string `json:"session_public_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
