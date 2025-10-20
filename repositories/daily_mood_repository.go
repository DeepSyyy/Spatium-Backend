package repositories

import (
	"fmt"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailyMoodRepository interface {
	CreateOrUpdate(userID int64, moodTag int64, note string, date time.Time) (*models.DailyMood, error)
	GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error)
	GetWeeklyMoodsByUserID(userID int64) ([]models.DailyMood, error)
	GetMoodStatistics(userID int64, days int) ([]models.MoodStats, error)
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

func (r *dailyMoodRepository) GetMoodStatistics(userID int64, days int) ([]models.MoodStats, error) {
	var stats []models.MoodStats

	query := fmt.Sprintf(`
		SELECT 
			dm.date,
			dm.mood_tag_internal_id AS mood_tag_id,
			mt.mood_name AS mood_tag_label,
			COUNT(dm.mood_tag_internal_id) AS count
		FROM daily_moods dm
		JOIN mood_tags mt ON dm.mood_tag_internal_id = mt.internal_id
		WHERE dm.user_internal_id = ?
		AND dm.date >= CURRENT_DATE - INTERVAL '%d day'
		GROUP BY dm.date, dm.mood_tag_internal_id, mt.mood_name
		ORDER BY dm.date ASC;
	`, days)

	err := r.db.Raw(query, userID).Scan(&stats).Error
	return stats, err
}
