package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}
