package seed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateRecoveryCode() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, 12)
	for i := range code {
		if i > 0 && i%4 == 0 {
			code[i] = '-'
		} else {
			code[i] = charset[rand.Intn(len(charset))]
		}
	}
	return string(code)
}

func SeedUsers(db *gorm.DB) {
	users := []models.User{
		{
			PublicID:     uuid.New(),
			Alias:        "Anon-0001",
			RecoveryCode: GenerateRecoveryCode(),
			CreatedAt:    time.Now(),
		},
		{
			PublicID:     uuid.New(),
			Alias:        "Anon-0002",
			RecoveryCode: GenerateRecoveryCode(),
			CreatedAt:    time.Now(),
		},
		{
			PublicID:     uuid.New(),
			Alias:        "Anon-0003",
			RecoveryCode: GenerateRecoveryCode(),
			CreatedAt:    time.Now(),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("❌ Failed to seed user %s: %v\n", user.Alias, err)
		} else {
			fmt.Printf("✅ Seeded user: %s (%s)\n", user.Alias, user.PublicID)
		}
	}
}
