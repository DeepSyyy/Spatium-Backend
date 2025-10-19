package controllers

import (
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

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	mood, err := c.service.CreateOrUpdate(userID, req.MoodTagID, req.Note)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to save daily mood", err.Error())
	}

	resp := models.DailyMoodResponse{
		PublicID:  mood.PublicID.String(),
		MoodTagID: mood.MoodTagInternalID,
		Note:      mood.Note,
		Date:      mood.Date.Format("2006-01-02"),
		CreatedAt: mood.CreatedAt,
		UpdatedAt: mood.UpdatedAt,
	}

	return utils.Success(ctx, "Daily mood saved successfully", resp)
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
