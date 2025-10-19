package controllers

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ChatSessionController struct {
	service services.ChatSessionService
}

func NewChatSessionController(service services.ChatSessionService) *ChatSessionController {
	return &ChatSessionController{service}
}

func (c *ChatSessionController) CreateSession(ctx *fiber.Ctx) error {
	var req struct {
		Title string `json:"title"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	session, err := c.service.Create(userID, req.Title)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to create session", err.Error())
	}

	return utils.Success(ctx, "Session created successfully", models.ChatSessionResponse{
		PublicID:  session.PublicID.String(),
		Title:     session.Title,
		MoodTag:   session.MoodTag,
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
	})
}

func (c *ChatSessionController) GetUserSessions(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	sessions, err := c.service.GetUserSessions(userID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to retrieve sessions", err.Error())
	}

	sessionsResponse := make([]models.ChatSessionResponse, len(sessions))
	for i, session := range sessions {
		sessionsResponse[i] = models.ChatSessionResponse{
			PublicID:  session.PublicID.String(),
			Title:     session.Title,
			MoodTag:   session.MoodTag,
			CreatedAt: session.CreatedAt.Format(time.RFC3339),
			UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
		}
	}

	return utils.Success(ctx, "Sessions retrieved successfully", sessionsResponse)
}

func (c *ChatSessionController) DeleteSession(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if publicID == "" {
		return utils.BadRequest(ctx, "Session ID is required", "")
	}

	err := c.service.DeleteSession(publicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to delete session", err.Error())
	}

	return utils.Success(ctx, "Session deleted successfully", nil)
}
