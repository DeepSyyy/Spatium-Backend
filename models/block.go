package models

import (
	"time"

	"github.com/google/uuid"
)

// Block represents a user blocking another user
type Block struct {
	InternalID    int64     `json:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID      uuid.UUID `json:"public_id" gorm:"uniqueIndex;not null"`
	BlockerID     int64     `json:"blocker_id" gorm:"not null;index"`      // User who is blocking
	BlockedUserID int64     `json:"blocked_user_id" gorm:"not null;index"` // User being blocked
	Reason        string    `json:"reason" gorm:"type:text"`               // Optional reason
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// BlockRequest is the request body for blocking a user
type BlockRequest struct {
	UserID string `json:"user_id" validate:"required"` // Public ID of user to block
	Reason string `json:"reason"`
}

// BlockResponse is the response for a block action
type BlockResponse struct {
	PublicID         string    `json:"public_id"`
	BlockedUserID    string    `json:"blocked_user_id"`
	BlockedUserAlias string    `json:"blocked_user_alias,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// BlockedUserInfo contains information about a blocked user
type BlockedUserInfo struct {
	PublicID  string    `json:"public_id"`
	Alias     string    `json:"alias"`
	BlockedAt time.Time `json:"blocked_at"`
}
