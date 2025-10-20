package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type AIReflectionRepository interface {
	Create(reflection *models.AIReflection) error
	GetByUser(userID int64) ([]models.AIReflection, error)
}

type aiReflectionRepository struct {
	db *gorm.DB
}

func NewAIReflectionRepository(db *gorm.DB) AIReflectionRepository {
	return &aiReflectionRepository{db}
}

func (r *aiReflectionRepository) Create(reflection *models.AIReflection) error {
	return r.db.Create(reflection).Error
}

func (r *aiReflectionRepository) GetByUser(userID int64) ([]models.AIReflection, error) {
	var reflections []models.AIReflection
	err := r.db.Where("user_internal_id = ?", userID).Order("created_at desc").Find(&reflections).Error
	return reflections, err
}
