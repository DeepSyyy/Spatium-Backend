package repositories

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByRecoveryCode(code string) (*models.User, error)
	FindByGoogleID(googleID string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByPublicID(publicID interface{}) (*models.User, error)
	UpdateLastLogin(user *models.User) error
	Update(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByPublicID(publicID string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByRecoveryCode(code string) (*models.User, error) {
	var user models.User
	err := r.db.Where("recovery_code = ?", code).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLogin(user *models.User) error {
	return r.db.Model(&models.User{}).
		Where("internal_id = ?", user.InternalID).
		Updates(map[string]interface{}{
			"last_login": time.Now(),
		}).Error
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("internal_id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("public_id = ?", publicID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPublicID(publicID interface{}) (*models.User, error) {
	var user models.User
	err := r.db.Where("public_id = ?", publicID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
