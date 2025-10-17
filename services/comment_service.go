package services

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
)

type CommentService interface {
	Create(userID int64, postID int64, content string) (*models.Comment, error)
	GetCommentsByPostID(postID int64) ([]models.Comment, error)
	Delete(commentID int64) error
	GetCommentInternalIDByPublicID(publicID string) (int64, error)
	GetCommentsByIDs(commentIDs []int64) ([]models.Comment, error)
}

type commentService struct {
	repo repositories.CommentRepository
}

func NewCommentService(repo repositories.CommentRepository) CommentService {
	return &commentService{repo}
}

func (s *commentService) Create(userID int64, postID int64, content string) (*models.Comment, error) {
	comment := &models.Comment{
		PublicID:  uuid.New(),
		UserID:    userID,
		PostID:    postID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) GetCommentsByPostID(postID int64) ([]models.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}

func (s *commentService) Delete(commentID int64) error {
	return s.repo.Delete(commentID)
}

func (s *commentService) GetCommentsByIDs(commentIDs []int64) ([]models.Comment, error) {
	return s.repo.GetCommentsByIDs(commentIDs)
}

func (s *commentService) GetCommentInternalIDByPublicID(publicID string) (int64, error) {
	return s.repo.GetCommentInternalIDByPublicID(publicID)
}
