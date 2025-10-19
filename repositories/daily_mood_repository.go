package repositories

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailyMoodRepository interface {
	CreateOrUpdate(userID int64, moodTag int64, note string, date time.Time) (*models.DailyMood, error)
	GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error)
	GetWeeklyMoodsByUserID(userID int64) ([]models.DailyMood, error)
}

type dailyMoodRepository struct {
	db *gorm.DB
}

func NewDailyMoodRepository(db *gorm.DB) DailyMoodRepository {
	return &dailyMoodRepository{db}
}

func (r *dailyMoodRepository) CreateOrUpdate(userID int64, moodTagID int64, note string, date time.Time) (*models.DailyMood, error) {
	var mood models.DailyMood
	err := r.db.Where("user_internal_id = ? AND date = ?", userID, date).First(&mood).Error

	if err == gorm.ErrRecordNotFound {
		mood = models.DailyMood{
			PublicID:          uuid.New(),
			UserInternalID:    userID,
			MoodTagInternalID: moodTagID,
			Note:              note,
			Date:              date,
		}
		if err := r.db.Create(&mood).Error; err != nil {
			return nil, err
		}
		return &mood, nil
	}

	if err != nil {
		return nil, err
	}

	// update jika sudah ada
	mood.MoodTagInternalID = moodTagID
	mood.Note = note
	mood.UpdatedAt = time.Now()
	if err := r.db.Save(&mood).Error; err != nil {
		return nil, err
	}

	return &mood, nil
}

func (r *dailyMoodRepository) GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error) {
	var mood models.DailyMood
	err := r.db.Where("user_internal_id = ? AND date = ?", userID, date).First(&mood).Error
	if err != nil {
		return nil, err
	}
	return &mood, nil
}

func (r *dailyMoodRepository) GetWeeklyMoodsByUserID(userID int64) ([]models.DailyMood, error) {
	var moods []models.DailyMood
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	err := r.db.Where("user_internal_id = ? AND date >= ?", userID, oneWeekAgo).Order("date ASC").Find(&moods).Error
	return moods, err
}
