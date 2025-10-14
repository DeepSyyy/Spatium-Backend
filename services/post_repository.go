package services

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
)

type PostService interface {
	Create(userID int64, moodTag string, content string) (*models.Post, error)
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) PostService {
	return &postService{repo}
}

func (s *postService) Create(userID int64, moodTag string, content string) (*models.Post, error) {
	post := &models.Post{
		PublicID:   uuid.New(),
		UserID:     userID,
		MoodTag:    moodTag,
		Content:    content,
		AiResponse: "",
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	return post, nil
}
