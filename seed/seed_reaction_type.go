package seed

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

func SeedReactionTypes(db *gorm.DB) {
	var reactions = []models.ReactionType{
		{Name: "neutral", Emoji: "😐", Description: "No specific emotion or reaction"},
		{Name: "like", Emoji: "👍", Description: "Positive or agreement reaction"},
		{Name: "love", Emoji: "❤️", Description: "Feeling empathy or strong positive emotion"},
		{Name: "laugh", Emoji: "😂", Description: "Something funny or joyful"},
		{Name: "sad", Emoji: "😢", Description: "Feeling empathy or sorrow"},
		{Name: "angry", Emoji: "😡", Description: "Feeling frustration or disagreement"},
		{Name: "tired", Emoji: "🥱", Description: "Feeling exhaustion or drained"},
	}

	for _, reaction := range reactions {
		var existing models.ReactionType
		err := config.DB.Where("name = ?", reaction.Name).First(&existing).Error
		if err == nil {
			continue // sudah ada, skip
		}
		if err := config.DB.Create(&reaction).Error; err != nil {
			log.Printf("❌ Failed to seed reaction: %s — %v\n", reaction.Name, err)
		} else {
			log.Printf("✅ Seeded reaction: %s (%s)\n", reaction.Name, reaction.Emoji)
		}
	}
}
