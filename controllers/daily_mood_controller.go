package controllers

import (
	"fmt"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type DailyMoodController struct {
	service services.DailyMoodService
}

func NewDailyMoodController(service services.DailyMoodService) *DailyMoodController {
	return &DailyMoodController{service}
}

func (c *DailyMoodController) CreateOrUpdate(ctx *fiber.Ctx) error {
	var req struct {
		MoodTagID int64  `json:"mood_tag_internal_id"`
		Note      string `json:"note"`
	}

	// 🔹 Parsing request body
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	// 🔹 Validasi token (ambil user_id dari JWT)
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	// 🔹 Simpan mood dan generate refleksi AI otomatis
	mood, reflection, err := c.service.CreateOrUpdate(userID, req.MoodTagID, req.Note)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to save daily mood", err.Error())
	}

	// 🔹 Format response ke frontend
	resp := fiber.Map{
		"mood": fiber.Map{
			"mood_tag_internal_id": req.MoodTagID,
			"note":                 req.Note,
			"date":                 mood.Date,
			"created_at":           mood.CreatedAt,
			"updated_at":           mood.UpdatedAt,
		},
	}

	// 🔹 Tambahkan refleksi AI kalau berhasil
	if reflection != nil {
		resp["reflection"] = fiber.Map{
			"mood_tag_internal_id": reflection.MoodTagID,
			"reflection":           reflection.Reflection,
			"created_at":           reflection.CreatedAt,
		}
	}

	return utils.Success(ctx, "Daily mood and reflection saved successfully", resp)
}

func (c *DailyMoodController) GetTodayMood(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	today := time.Now().Truncate(24 * time.Hour)
	mood, err := c.service.GetByUserAndDate(userID, today)
	if err != nil {
		return utils.BadRequest(ctx, "No mood entry found for today", err.Error())
	}

	resp := models.DailyMoodResponse{
		PublicID:  mood.PublicID.String(),
		MoodTagID: mood.MoodTagInternalID,
		Note:      mood.Note,
		Date:      mood.Date.Format("2006-01-02"),
		CreatedAt: mood.CreatedAt,
		UpdatedAt: mood.UpdatedAt,
	}

	return utils.Success(ctx, "Today's mood retrieved successfully", resp)
}

func (c *DailyMoodController) GetWeeklyMoods(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	moods, err := c.service.GetWeeklyMoods(userID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to retrieve weekly moods", err.Error())
	}

	resp := make([]models.DailyMoodResponse, len(moods))
	for i, m := range moods {
		resp[i] = models.DailyMoodResponse{
			PublicID:  m.PublicID.String(),
			MoodTagID: m.MoodTagInternalID,
			Note:      m.Note,
			Date:      m.Date.Format("2006-01-02"),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return utils.Success(ctx, "Weekly moods retrieved successfully", resp)
}

func (c *DailyMoodController) GetMoodStatistics(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	// ✅ Ambil query param ?days, default 7
	days := ctx.QueryInt("days", 7)

	// ✅ Ambil data statistik dari service
	stats, err := c.service.GetMoodStatistics(userID, days)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to fetch mood statistics", err.Error())
	}

	// ✅ Bentuk response array yang bersih
	response := make([]fiber.Map, len(stats))
	for i, s := range stats {
		response[i] = fiber.Map{
			"date":                 s.Date.Format("2006-01-02"),
			"mood_tag_internal_id": s.MoodTagID,
			"mood_tag_label":       s.MoodTagLabel,
			"count":                s.Count,
		}
	}

	return utils.Success(ctx, fmt.Sprintf("Mood statistics for last %d days retrieved successfully", days), response)
}

func (c *DailyMoodController) GetMoodChart(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	days := ctx.QueryInt("days", 7)
	if days < 1 {
		days = 7
	}

	chartData, err := c.service.GetMoodChart(userID, days)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to generate mood chart", err.Error())
	}

	return utils.Success(ctx, fmt.Sprintf("Mood chart for last %d days", days), chartData)
}
