package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type AIReflectionController struct {
	service services.AIReflectionService
}

func NewAIReflectionController(service services.AIReflectionService) *AIReflectionController {
	return &AIReflectionController{service}
}

func (c *AIReflectionController) GenerateReflection(ctx *fiber.Ctx) error {
	var req struct {
		MoodTag  int64  `json:"mood_internal_id"`
		MoodDesc string `json:"mood_description"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid input", err.Error())
	}

	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	reflection, err := c.service.GenerateReflection(userID, req.MoodTag, req.MoodDesc)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to generate AI reflection", err.Error())
	}

	resp := models.AIReflectionResponse{
		PublicID:   reflection.PublicID.String(),
		MoodTagID:  reflection.MoodTagID,
		Reflection: reflection.Reflection,
		CreatedAt:  reflection.CreatedAt,
	}

	return utils.Success(ctx, "AI reflection generated successfully", resp)
}

func (c *AIReflectionController) GetUserReflection(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	reflections, err := c.service.GetUserReflections(userID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to fetch reflections", err.Error())
	}

	resp := make([]models.AIReflectionResponse, len(reflections))
	for i, r := range reflections {
		resp[i] = models.AIReflectionResponse{
			PublicID:   r.PublicID.String(),
			MoodTagID:  r.MoodTagID,
			Reflection: r.Reflection,
			CreatedAt:  r.CreatedAt,
		}
	}

	return utils.Success(ctx, "Reflections retrieved successfully", resp)
}
