package models

import "github.com/google/uuid"

type User struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id" gorm:"uniqueIndex;not null"`
	CreatedAt  int64     `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	Name       string    `json:"name" db:"name" gorm:"type:varchar(100);not null"`
}
