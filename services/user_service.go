package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	Register() (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register() (*models.User, error) {
	rand.Seed(time.Now().UnixNano())
	user := &models.User{
		PublicID:     uuid.New(),
		Alias:        fmt.Sprintf("Anon-%04d", rand.Intn(9999)),
		RecoveryCode: generateRecoveryCode(),
		CreatedAt:    time.Now(),
	}

	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func generateRecoveryCode() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
