package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	commentService services.CommentService
	postService    services.PostService
}

func NewCommentController(commentService services.CommentService, postService services.PostService) *CommentController {
	return &CommentController{commentService: commentService, postService: postService}
}

func (c *CommentController) CreateComment(ctx *fiber.Ctx) error {
	var req struct {
		PostID  string `json:"post_id"`
		Content string `json:"content"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid input", err.Error())
	}

	// Ambil user ID dari token JWT (diset di middleware)
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}
	if req.Content == "" {
		return utils.BadRequest(ctx, "Content cannot be empty", "")
	}
	// Convert PostID from public ID (string) to internal ID (int64)
	postInternalID, err := c.postService.GetPostInternalIDByPublicID(req.PostID)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid Post ID", err.Error())
	}
	comment, err := c.commentService.Create(userID, postInternalID, req.Content)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to create comment", err.Error())
	}

	resp := models.CommentResponse{
		PublicID:  comment.PublicID.String(),
		PostID:    req.PostID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
	return utils.Success(ctx, "Comment created successfully", resp)

}

func (c *CommentController) GetCommentsByPostID(ctx *fiber.Ctx) error {
	postID := ctx.Params("post_id")
	if postID == "" {
		return utils.BadRequest(ctx, "Post ID is required", "")
	}

	// Convert PostID from public ID (string) to internal ID (int64)
	postInternalID, err := c.postService.GetPostInternalIDByPublicID(postID)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid Post ID", err.Error())
	}

	comments, err := c.commentService.GetCommentsByPostID(postInternalID)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to retrieve comments", err.Error())
	}

	var resp []models.CommentResponse
	for _, comment := range comments {
		resp = append(resp, models.CommentResponse{
			PublicID:  comment.PublicID.String(),
			PostID:    postID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		})
	}

	return utils.Success(ctx, "Comments retrieved successfully", resp)
}

func (c *CommentController) DeleteComment(ctx *fiber.Ctx) error {
	commentID := ctx.Params("comment_id")
	if commentID == "" {
		return utils.BadRequest(ctx, "Comment ID is required", "")
	}

	// Convert CommentID from public ID (string) to internal ID (int64)
	commentInternalID, err := c.commentService.GetCommentInternalIDByPublicID(commentID)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid Comment ID", err.Error())
	}

	if err := c.commentService.Delete(commentInternalID); err != nil {
		return utils.BadRequest(ctx, "Failed to delete comment", err.Error())
	}

	return utils.Success(ctx, "Comment deleted successfully", nil)
}
