package models

type MoodTag struct {
	InternalID  int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	MoodName    string    `json:"mood_name" db:"mood_name" gorm:"type:varchar(50);unique;not null"`
	Description string    `json:"description" db:"description" gorm:"type:text"`
}