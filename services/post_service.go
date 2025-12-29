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

// PostCreateResult contains the created post and any additional information
type PostCreateResult struct {
	Post                    *models.Post
	ContainsCrisisIndicator bool
	CrisisSupportMessage    string
}

type PostService interface {
	Create(userID int64, moodTag int64, content string) (*PostCreateResult, error)
	GetAll() ([]models.Post, error)
	GetAllFiltered(excludeUserIDs []int64) ([]models.Post, error)
	GetPostDetail(publicID string) (*models.Post, error)
	GetPostsByUserID(userID int64) ([]models.Post, error)
	GetPostInternalIDByPublicID(publicID string) (int64, error)
	Update(publicID string, updatedPost *models.Post) error
	Delete(publicID string) error
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo}
}

func (s *postService) Create(userID int64, moodTag int64, content string) (*PostCreateResult, error) {
	result := &PostCreateResult{}

	// Step 1: Content Moderation - Check for harmful content
	moderationResult, err := utils.ModerateContent(content)
	if err != nil {
		log.Println("Warning: Content moderation failed:", err)
		// Continue with post creation even if moderation fails
	} else if moderationResult != nil && moderationResult.IsBlocked {
		// Content is blocked - return error with reason
		return nil, errors.New("Konten tidak dapat dipublikasikan: " + moderationResult.Reason)
	}

	// Step 2: Check for crisis indicators (for supportive response)
	if utils.ContainsCrisisIndicators(content) {
		result.ContainsCrisisIndicator = true
		result.CrisisSupportMessage = utils.GetCrisisSupportMessage()
	}

	// Step 3: Generate AI empathic response
	aiResp, err := utils.GenerateEmpathicResponse(content, int(moodTag))
	if err != nil {
		log.Println("Error generating AI response:", err)
		return nil, err
	}

	// Step 4: Create the post
	post := &models.Post{
		PublicID:   uuid.New(),
		UserID:     userID,
		MoodTagID:  moodTag,
		Content:    content,
		AiResponse: aiResp,
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	result.Post = post
	return result, nil
}

func (s *postService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}

// GetAllFiltered returns posts excluding those from specified user IDs (blocked users)
func (s *postService) GetAllFiltered(excludeUserIDs []int64) ([]models.Post, error) {
	if len(excludeUserIDs) == 0 {
		return s.repo.GetAll()
	}
	return s.repo.GetAllExcluding(excludeUserIDs)
}

func (s *postService) GetPostDetail(publicID string) (*models.Post, error) {
	return s.repo.GetPostDetail(publicID)
}

func (s *postService) GetPostsByUserID(userID int64) ([]models.Post, error) {
	return s.repo.GetPostsByUserID(userID)
}

func (s *postService) Update(publicID string, updatedPost *models.Post) error {
	return s.repo.Update(publicID, updatedPost)
}

func (s *postService) Delete(publicID string) error {
	return s.repo.Delete(publicID)
}

func (s *postService) GetPostInternalIDByPublicID(publicID string) (int64, error) {
	return s.repo.GetPostInternalIDByPublicID(publicID)
}
