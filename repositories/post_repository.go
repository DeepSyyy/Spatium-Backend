package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetAll() ([]models.Post, error)
	GetPostDetail(publicID string) (*models.Post, error)
	GetPostsByUserID(userID int64) ([]models.Post, error)
	Update(publicID string, updatedPost *models.Post) error
	Delete(publicID string) error
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

func (r *postRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Find(&posts).Error
	return posts, err
}

func (r *postRepository) GetPostDetail(publicID string) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("public_id = ?", publicID).First(&post).Error
	return &post, err
}

func (r *postRepository) GetPostsByUserID(userID int64) ([]models.Post, error) {
	var posts []models.Post
	err := config.DB.Where("user_internal_id = ?", userID).Find(&posts).Order("created_at desc").Error
	return posts, err
}

func (r *postRepository) Update(publicID string, updatedPost *models.Post) error {
	return r.db.Model(&models.Post{}).Where("public_id = ?", publicID).Updates(updatedPost).Error
}

func (r *postRepository) Delete(publicID string) error {
	return r.db.Where("public_id = ?", publicID).Delete(&models.Post{}).Error
}
