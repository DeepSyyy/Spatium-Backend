package seed

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

func SeedMoodTags(db *gorm.DB) {
	defaultMoods := []models.MoodTag{
		{MoodName: "neutral", Description: "Feeling balanced or indifferent"},
		{MoodName: "happy", Description: "Feeling positive or cheerful"},
		{MoodName: "sad", Description: "Feeling down or upset"},
		{MoodName: "anxious", Description: "Feeling nervous or worried"},
		{MoodName: "angry", Description: "Feeling irritated or frustrated"},
		{MoodName: "tired", Description: "Feeling physically or mentally drained"},
	}

	for _, mood := range defaultMoods {
		var existing models.MoodTag
		err := db.Where("mood_name = ?", mood.MoodName).First(&existing).Error
		if err == nil {
			continue // already exists, skip
		}
		if err := db.Create(&mood).Error; err != nil {
			println("❌ Failed to seed mood tag:", mood.MoodName, "—", err.Error())
		} else {
			println("✅ Seeded mood tag:", mood.MoodName)
		}
	}
}
