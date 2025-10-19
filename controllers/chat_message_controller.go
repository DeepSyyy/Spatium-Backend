package controllers

import (
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ChatMessageController struct {
	service services.ChatMessageService
}

func NewChatMessageController(service services.ChatMessageService) *ChatMessageController {
	return &ChatMessageController{service}
}

func (c *ChatMessageController) Create(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	sessionPublicID := ctx.Params("session_id")
	if sessionPublicID == "" {
		return utils.BadRequest(ctx, "Session ID is required", "")
	}

	// ✅ Validate ownership first
	sessionIDint, err := c.service.ValidateSessionOwnership(userID, sessionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Unauthorized", err.Error())
	}

	var req struct {
		Content string `json:"content"`
		Sender  string `json:"sender"` // optional, default = user
	}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}
	if req.Sender == "" {
		req.Sender = "user"
	}

	message, err := c.service.Create(req.Content, sessionIDint, req.Sender)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to create message", err.Error())
	}

	resp := models.ChatMessageResponse{
		SessionID: sessionPublicID,
		Sender:    message.Sender,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
	}

	return utils.Success(ctx, "Message created successfully", resp)
}

func (c *ChatMessageController) GetMessagesBySession(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	sessionPublicID := ctx.Params("session_id")
	sessionIDint, err := c.service.ValidateSessionOwnership(userID, sessionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Unauthorized", err.Error())
	}

	limit := ctx.QueryInt("limit", 20)
	messages, err := c.service.FindBySession(sessionIDint, limit)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to fetch messages", err.Error())
	}

	resp := make([]models.ChatMessageResponse, len(messages))
	for i, m := range messages {
		resp[i] = models.ChatMessageResponse{
			SessionID: sessionPublicID,
			Sender:    m.Sender,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		}
	}
	return utils.Success(ctx, "Messages retrieved successfully", resp)
}

func (c *ChatMessageController) GetLastMessages(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	sessionPublicID := ctx.Params("session_id")
	sessionIDint, err := c.service.ValidateSessionOwnership(userID, sessionPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Unauthorized", err.Error())
	}

	limit := ctx.QueryInt("limit", 10)
	messages, err := c.service.FindLastMessages(sessionIDint, limit)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to fetch last messages", err.Error())
	}

	resp := make([]models.ChatMessageResponse, len(messages))
	for i, m := range messages {
		resp[i] = models.ChatMessageResponse{
			SessionID: sessionPublicID,
			Sender:    m.Sender,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
		}
	}
	return utils.Success(ctx, "Last messages retrieved successfully", resp)
}
