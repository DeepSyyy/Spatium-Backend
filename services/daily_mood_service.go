package services

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
)

type DailyMoodService interface {
	CreateOrUpdate(userID, moodTagID int64, note string) (*models.DailyMood, error)
	GetWeeklyMoods(userID int64) ([]models.DailyMood, error)
	GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error)
}

type dailyMoodService struct {
	repo repositories.DailyMoodRepository
}

func NewDailyMoodService(repo repositories.DailyMoodRepository) DailyMoodService {
	return &dailyMoodService{repo}
}

func (s *dailyMoodService) CreateOrUpdate(userID, moodTagID int64, note string) (*models.DailyMood, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.repo.CreateOrUpdate(userID, moodTagID, note, today)
}

func (s *dailyMoodService) GetWeeklyMoods(userID int64) ([]models.DailyMood, error) {
	return s.repo.GetWeeklyMoodsByUserID(userID)
}

func (s *dailyMoodService) GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error) {
	return s.repo.GetByUserAndDate(userID, date)
}
