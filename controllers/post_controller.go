package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
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
		req.MoodTag = 1 // default mood
	}

	post, err := c.service.Create(userID, req.MoodTag, req.Content)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to create post", err.Error())
	}

	resp := models.PostResponse{
		PublicID:   post.PublicID.String(),
		Content:    post.Content,
		AiResponse: post.AiResponse,
		MoodTagID:  post.MoodTagID,
		CreatedAt:  post.CreatedAt,
	}

	return utils.Success(ctx, "Post created successfully", resp)
}

// GetAllPosts handles GET /posts
func (c *PostController) GetAllPosts(ctx *fiber.Ctx) error {
	posts, err := c.service.GetAll()
	if err != nil {
		return utils.BadRequest(ctx, "Failed to retrieve posts", err.Error())
	}
	// Map ke response tanpa internal_id
	resp := make([]models.PostResponse, len(posts))

	for i, post := range posts {
		resp[i] = models.PostResponse{
			PublicID:   post.PublicID.String(),
			Content:    post.Content,
			AiResponse: post.AiResponse,
			MoodTagID:  post.MoodTagID,
			CreatedAt:  post.CreatedAt,
		}
	}

	return utils.Success(ctx, "Posts retrieved successfully", resp)
}

// GetPostDetail handles GET /posts/:public_id
func (c *PostController) GetPostDetail(ctx *fiber.Ctx) error {
	publicID := ctx.Params("post_id")
	if publicID == "" {
		return utils.BadRequest(ctx, "Post ID is required", "")
	}

	post, err := c.service.GetPostDetail(publicID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to retrieve post", err.Error())
	}

	// Map ke response tanpa internal_id
	resp := models.PostResponse{
		PublicID:   post.PublicID.String(),
		Content:    post.Content,
		AiResponse: post.AiResponse,
		MoodTagID:  post.MoodTagID,
		CreatedAt:  post.CreatedAt,
	}

	return utils.Success(ctx, "Post retrieved successfully", resp)
}

// GetPostsByUser handles GET /users/:user_id/posts
func (c *PostController) GetPostsByUser(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	posts, err := c.service.GetPostsByUserID(userID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to retrieve user's posts", err.Error())
	}

	resp := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		resp[i] = models.PostResponse{
			PublicID:   post.PublicID.String(),
			Content:    post.Content,
			AiResponse: post.AiResponse,
			MoodTagID:  post.MoodTagID,
			CreatedAt:  post.CreatedAt,
		}
	}

	return utils.Success(ctx, "User's posts retrieved successfully", resp)
}

// UpdatePost handles PUT /posts/:public_id
func (c *PostController) UpdatePost(ctx *fiber.Ctx) error {
	publicID := ctx.Params("post_id")
	if publicID == "" {
		return utils.BadRequest(ctx, "Post ID is required", "")
	}

	var req struct {
		Content string `json:"content"`
		MoodTag int64  `json:"mood_internal_id"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid input", err.Error())
	}

	updatedPost := &models.Post{
		Content:   req.Content,
		MoodTagID: req.MoodTag,
	}

	if err := c.service.Update(publicID, updatedPost); err != nil {
		return utils.BadRequest(ctx, "Failed to update post", err.Error())
	}

	resp := models.PostResponse{
		PublicID:   publicID,
		Content:    updatedPost.Content,
		MoodTagID:  updatedPost.MoodTagID,
		AiResponse: updatedPost.AiResponse,
	}

	return utils.Success(ctx, "Post updated successfully", resp)
}

// DeletePost handles DELETE /posts/:public_id
func (c *PostController) DeletePost(ctx *fiber.Ctx) error {
	publicID := ctx.Params("post_id")
	if publicID == "" {
		return utils.BadRequest(ctx, "Post ID is required", "")
	}

	if err := c.service.Delete(publicID); err != nil {
		return utils.BadRequest(ctx, "Failed to delete post", err.Error())
	}

	return utils.Success(ctx, "Post deleted successfully", nil)
}
