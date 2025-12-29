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
	Register(alias string) (*models.User, error)
	Login(recoveryCode string) (*models.User, error)
	GoogleLogin(googleID, email, alias, photoURL string) (*models.User, error)
	UpdateAlias(userID, alias string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(alias string) (*models.User, error) {
	rand.Seed(time.Now().UnixNano())

	// Use provided alias or generate anonymous one
	if alias == "" {
		alias = fmt.Sprintf("Anon-%04d", rand.Intn(9999))
	}

	user := &models.User{
		PublicID:     uuid.New(),
		Alias:        alias,
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

func (s *userService) Login(recoveryCode string) (*models.User, error) {
	user, err := s.repo.FindByRecoveryCode(recoveryCode)
	if err != nil {
		return nil, err
	}

	// Update last login time
	user.LastLogin = time.Now()
	err = s.repo.UpdateLastLogin(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GoogleLogin(googleID, email, alias, photoURL string) (*models.User, error) {
	// Try to find existing user by Google ID
	user, err := s.repo.FindByGoogleID(googleID)
	if err == nil {
		// User exists, update last login
		user.LastLogin = time.Now()
		_ = s.repo.UpdateLastLogin(user)
		return user, nil
	}

	// User doesn't exist, create new one
	if alias == "" {
		// Extract name from email or use default
		alias = email
		if len(email) > 0 {
			if idx := len(email); idx > 20 {
				alias = email[:20]
			}
		}
	}

	user = &models.User{
		PublicID:     uuid.New(),
		GoogleID:     googleID,
		Email:        email,
		Alias:        alias,
		PhotoURL:     photoURL,
		RecoveryCode: generateRecoveryCode(), // Generate unique recovery code
		CreatedAt:    time.Now(),
		LastLogin:    time.Now(),
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateAlias(userID, alias string) (*models.User, error) {
	// Find user by public ID
	publicID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update alias
	user.Alias = alias
	err = s.repo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update alias: %w", err)
	}

	return user, nil
}
