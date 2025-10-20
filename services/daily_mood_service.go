package services

import (
	"sort"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
)

type DailyMoodService interface {
	CreateOrUpdate(userID, moodTagID int64, note string) (*models.DailyMood, *models.AIReflection, error)
	GetWeeklyMoods(userID int64) ([]models.DailyMood, error)
	GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error)
	GetMoodStatistics(userID int64, days int) ([]models.MoodStats, error)
	GetMoodChart(userID int64, days int) (map[string]interface{}, error)
}

type dailyMoodService struct {
	repo         repositories.DailyMoodRepository
	aiReflection AIReflectionService
}

func NewDailyMoodService(repo repositories.DailyMoodRepository, aiReflection AIReflectionService) DailyMoodService {
	return &dailyMoodService{repo, aiReflection}
}

func (s *dailyMoodService) CreateOrUpdate(userID, moodTagID int64, note string) (*models.DailyMood, *models.AIReflection, error) {
	// 🔹 Pastikan hanya ada 1 mood per hari (truncate waktu ke hari ini)
	today := time.Now().Truncate(24 * time.Hour)

	// 1️⃣ Simpan atau update daily mood ke DB
	mood := &models.DailyMood{
		UserInternalID:    userID,
		MoodTagInternalID: moodTagID,
		Note:              note,
		Date:              today, // Use the time.Time value directly
	}

	mood, err := s.repo.CreateOrUpdate(mood.UserInternalID, mood.MoodTagInternalID, mood.Note, today)
	if err != nil {
		return nil, nil, err
	}

	// 2️⃣ Deskripsi mood untuk AI Reflection
	moodDescMap := map[int64]string{
		1: "neutral - feeling balanced or indifferent",
		2: "happy - feeling positive or cheerful",
		3: "sad - feeling down or upset",
		4: "anxious - feeling nervous or worried",
		5: "angry - feeling irritated or frustrated",
		6: "tired - feeling physically or mentally drained",
	}

	moodDesc, ok := moodDescMap[moodTagID]
	if !ok {
		moodDesc = "neutral - feeling balanced or indifferent"
	}

	// 3️⃣ Panggil AI Reflection Service
	reflection, err := s.aiReflection.GenerateReflection(userID, moodTagID, moodDesc)
	if err != nil {
		// Jika AI gagal, tetap kembalikan mood (biar user experience-nya gak stop)
		return mood, nil, err
	}

	return mood, reflection, nil
}

func (s *dailyMoodService) GetWeeklyMoods(userID int64) ([]models.DailyMood, error) {
	return s.repo.GetWeeklyMoodsByUserID(userID)
}

func (s *dailyMoodService) GetByUserAndDate(userID int64, date time.Time) (*models.DailyMood, error) {
	return s.repo.GetByUserAndDate(userID, date)
}

func (s *dailyMoodService) GetMoodStatistics(userID int64, days int) ([]models.MoodStats, error) {
	if days <= 0 {
		days = 7
	} else if days > 365 {
		days = 365
	}

	return s.repo.GetMoodStatistics(userID, days)
}

func (s *dailyMoodService) GetMoodChart(userID int64, days int) (map[string]interface{}, error) {
	stats, err := s.repo.GetMoodStatistics(userID, days)
	if err != nil {
		return nil, err
	}

	// Kumpulan tanggal unik
	dateMap := make(map[string]bool)
	for _, stat := range stats {
		dateMap[stat.Date.Format("2006-01-02")] = true
	}

	// Urutkan tanggalnya
	var labels []string
	for date := range dateMap {
		labels = append(labels, date)
	}
	sort.Strings(labels)

	// Siapkan semua kategori mood (agar tampil meskipun count=0)
	moodLabels := []string{"neutral", "happy", "sad", "anxious", "angry", "tired"}
	data := make(map[string][]int)
	for _, mood := range moodLabels {
		data[mood] = make([]int, len(labels))
	}

	// Isi data sesuai tanggal & mood
	for _, stat := range stats {
		dateStr := stat.Date.Format("2006-01-02")
		for i, label := range labels {
			if dateStr == label {
				data[stat.MoodTagLabel][i] = int(stat.Count) // Convert int64 to int
			}
		}
	}

	// Return format untuk frontend
	return map[string]interface{}{
		"labels": labels,
		"data":   data,
	}, nil
}
