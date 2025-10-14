package models

import "github.com/google/uuid"

type Reaction struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id"`
	UserID     int64     `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id;not null"`
	Emoji      string    `json:"emoji" db:"emoji" gorm:"type:varchar(10);not null"`
	CreatedAt  int64     `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
}
