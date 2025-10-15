package services

import (
	"log"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/google/uuid"
)

type PostService interface {
	Create(userID int64, moodTag int64, content string) (*models.Post, error)
	GetAll() ([]models.Post, error)
	GetPostDetail(publicID string) (*models.Post, error)
	GetPostsByUserID(userID int64) ([]models.Post, error)
	Update(publicID string, updatedPost *models.Post) error
	Delete(publicID string) error
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo}
}

func (s *postService) Create(userID int64, moodTag int64, content string) (*models.Post, error) {
	aiResp, err := utils.GenerateEmpathicResponse(content, int(moodTag))
	if err != nil {
		log.Println("Error generating AI response:", err)
		return nil, err
	}
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

	return post, nil
}

func (s *postService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
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
