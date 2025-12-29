package services

import (
	"errors"
	"log"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/google/uuid"
)

// CommentCreateResult contains the created comment and any additional information
type CommentCreateResult struct {
	Comment                 *models.Comment
	ContainsCrisisIndicator bool
	CrisisSupportMessage    string
}

type CommentService interface {
	Create(userID int64, postID int64, content string) (*CommentCreateResult, error)
	GetCommentsByPostID(postID int64) ([]models.Comment, error)
	GetCommentsByPostIDFiltered(postID int64, excludeUserIDs []int64) ([]models.Comment, error)
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

func (s *commentService) Create(userID int64, postID int64, content string) (*CommentCreateResult, error) {
	result := &CommentCreateResult{}

	// Step 1: Content Moderation - Check for harmful content
	moderationResult, err := utils.ModerateContent(content)
	if err != nil {
		log.Println("Warning: Content moderation failed:", err)
		// Continue with comment creation even if moderation fails
	} else if moderationResult != nil && moderationResult.IsBlocked {
		// Content is blocked - return error with reason
		return nil, errors.New("Komentar tidak dapat dipublikasikan: " + moderationResult.Reason)
	}

	// Step 2: Check for crisis indicators (for supportive response)
	if utils.ContainsCrisisIndicators(content) {
		result.ContainsCrisisIndicator = true
		result.CrisisSupportMessage = utils.GetCrisisSupportMessage()
	}

	// Step 3: Create the comment
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

	result.Comment = comment
	return result, nil
}

func (s *commentService) GetCommentsByPostID(postID int64) ([]models.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}

// GetCommentsByPostIDFiltered returns comments excluding those from specified user IDs
func (s *commentService) GetCommentsByPostIDFiltered(postID int64, excludeUserIDs []int64) ([]models.Comment, error) {
	if len(excludeUserIDs) == 0 {
		return s.repo.GetCommentsByPostID(postID)
	}
	return s.repo.GetCommentsByPostIDExcluding(postID, excludeUserIDs)
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
