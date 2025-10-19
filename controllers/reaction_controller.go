package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ReactionController struct {
	reactionService services.ReactionService
	postService     services.PostService
}

func NewReactionController(r services.ReactionService, p services.PostService) *ReactionController {
	return &ReactionController{r, p}
}

func (c *ReactionController) ReactToPost(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok || userID == 0 {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	postPublicID := ctx.Params("post_id")
	if postPublicID == "" {
		return utils.BadRequest(ctx, "Post ID is required", "")
	}

	var req struct {
		Emoji string `json:"emoji"` // New field for custom emoji
	}
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid input", "reaction_type_id is required")
	}

	reactionTypeID, err := c.reactionService.GetReactionTypeIDByEmoji(req.Emoji)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid emoji", err.Error())
	}

	//translate pubilc id to internal id
	post, err := c.postService.GetPostInternalIDByPublicID(postPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Post not found", err.Error())
	}
	err = c.reactionService.ReactToPost(userID, post, reactionTypeID, req.Emoji)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to react to post", err.Error())
	}

	reactResponse := models.ReactionResponse{
		Emoji: req.Emoji,
	}

	return utils.Success(ctx, "Reaction recorded successfully", reactResponse)
}

func (c *ReactionController) GetReactionSummary(ctx *fiber.Ctx) error {
	postPublicID := ctx.Params("post_id")
	if postPublicID == "" {
		return utils.BadRequest(ctx, "Post ID is required", "")
	}

	//translate pubilc id to internal id
	post, err := c.postService.GetPostInternalIDByPublicID(postPublicID)
	if err != nil {
		return utils.BadRequest(ctx, "Post not found", err.Error())
	}

	summary, err := c.reactionService.GetReactionSummary(post)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to get reaction summary", err.Error())
	}

	return utils.Success(ctx, "Reaction summary retrieved successfully", summary)
}
