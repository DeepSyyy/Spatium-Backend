package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type BlockController struct {
	blockService services.BlockService
}

func NewBlockController(blockService services.BlockService) *BlockController {
	return &BlockController{blockService: blockService}
}

// BlockUser handles POST /users/block
// @Summary Block a user
// @Description Block another user to hide their content and prevent interactions
func (c *BlockController) BlockUser(ctx *fiber.Ctx) error {
	var req struct {
		UserID string `json:"user_id"`
		Reason string `json:"reason"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Input tidak valid", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	if req.UserID == "" {
		return utils.BadRequest(ctx, "User ID wajib diisi", "")
	}

	block, err := c.blockService.Block(userID, req.UserID, req.Reason)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal memblokir user", err.Error())
	}

	return utils.Success(ctx, "User berhasil diblokir", fiber.Map{
		"public_id":  block.PublicID.String(),
		"blocked_at": block.CreatedAt,
	})
}

// UnblockUser handles DELETE /users/block/:user_id
// @Summary Unblock a user
// @Description Unblock a previously blocked user
func (c *BlockController) UnblockUser(ctx *fiber.Ctx) error {
	blockedUserID := ctx.Params("user_id")
	if blockedUserID == "" {
		return utils.BadRequest(ctx, "User ID wajib diisi", "")
	}

	// Get user ID from JWT token
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	err := c.blockService.Unblock(userID, blockedUserID)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal membuka blokir user", err.Error())
	}

	return utils.Success(ctx, "User berhasil dibuka blokirnya", nil)
}

// GetBlockedUsers handles GET /users/blocked
// @Summary Get blocked users
// @Description Get list of all users blocked by the current user
func (c *BlockController) GetBlockedUsers(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	blockedUsers, err := c.blockService.GetBlockedUsers(userID)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil daftar user yang diblokir", err.Error())
	}

	if blockedUsers == nil {
		blockedUsers = []models.BlockedUserInfo{}
	}

	return utils.Success(ctx, "Daftar user yang diblokir berhasil diambil", blockedUsers)
}
