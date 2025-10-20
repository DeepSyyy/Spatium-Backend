package models

import "time"

type MoodStats struct {
	Date         time.Time `json:"date"`
	MoodTagID    int64     `json:"mood_tag_internal_id"`
	MoodTagLabel string    `json:"mood_tag_label"`
	Count        int64     `json:"count"`
}
