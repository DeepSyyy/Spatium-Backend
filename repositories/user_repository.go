package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByRecoveryCode(code string) (*models.User, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByRecoveryCode(code string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("recovery_code = ?", code).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
