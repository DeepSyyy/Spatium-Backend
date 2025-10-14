package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type PostController struct {
	service services.PostService
}

func NewPostController(service services.PostService) *PostController {
	return &PostController{service}
}

// CreatePost handles POST /posts
func (c *PostController) CreatePost(ctx *fiber.Ctx) error {
	var req struct {
		Content string `json:"content"`
		MoodTag int64  `json:"mood_internal_id"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid input", err.Error())
	}

	// Ambil user ID dari token JWT (diset di middleware)
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	// Validasi input minimal
	if req.Content == "" {
		return utils.BadRequest(ctx, "Content cannot be empty", "")
	}
	//if moodtag == nil maka set default ke 0
	if req.MoodTag == 0 {
		req.MoodTag = 0 // default mood
	}

	post, err := c.service.Create(userID, req.MoodTag, req.Content)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to create post", err.Error())
	}

	return utils.Success(ctx, "Post created successfully", post)
}
