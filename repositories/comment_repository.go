package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetCommentsByPostID(postID int64) ([]models.Comment, error)
	GetCommentsByPostIDExcluding(postID int64, excludeUserIDs []int64) ([]models.Comment, error)
	Delete(commentID int64) error
	GetCommentInternalIDByPublicID(publicID string) (int64, error)
	GetCommentsByIDs(commentIDs []int64) ([]models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetCommentsByPostID(postID int64) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("post_internal_id = ?", postID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *commentRepository) GetCommentsByPostIDExcluding(postID int64, excludeUserIDs []int64) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("post_internal_id = ? AND user_internal_id NOT IN ?", postID, excludeUserIDs).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *commentRepository) Delete(commentID int64) error {
	return r.db.Where("internal_id = ?", commentID).Delete(&models.Comment{}).Error
}

func (r *commentRepository) GetCommentInternalIDByPublicID(publicID string) (int64, error) {
	var comment models.Comment
	err := r.db.Select("internal_id").Where("public_id = ?", publicID).First(&comment).Error
	if err != nil {
		return 0, err
	}
	return comment.InternalID, nil
}

func (r *commentRepository) GetCommentsByIDs(commentIDs []int64) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Where("internal_id IN ?", commentIDs).Find(&comments).Error
	return comments, err
}
