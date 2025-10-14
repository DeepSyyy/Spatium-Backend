package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	InternalID   int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID     uuid.UUID `json:"public_id" db:"public_id" gorm:"uniqueIndex;not null"`
	Alias        string    `gorm:"type:varchar(50);not null"`
	RecoveryCode string    `gorm:"type:char(12);uniqueIndex;not null"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	LastLogin    time.Time `gorm:"index"`
}

type UserResponse struct {
	PublicID  uuid.UUID  `json:"public_id"`
	Alias     string     `json:"alias"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}
