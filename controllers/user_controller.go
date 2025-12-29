package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{service: userService}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var body struct {
		Alias string `json:"alias"`
	}
	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	user, err := c.service.Register(body.Alias)
	if err != nil {
		// Handle error (e.g., return HTTP 500)
		return utils.BadRequest(ctx, "Registration failed", err.Error())
	}

	token, err := utils.GenerateToken(user.PublicID.String(), user.Alias)
	if err != nil {
		return utils.BadRequest(ctx, "Token generation failed", err.Error())
	}

	// Copy model → response DTO
	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)

	// Tambahkan token ke response
	response := fiber.Map{
		"user":  userResponse,
		"token": token,
	}

	return utils.Success(ctx, "Registration successful", response)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		RecoveryCode string `json:"recovery_code"`
	}
	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	user, err := c.service.Login(body.RecoveryCode)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Failed, invalid credential", err.Error())
	}

	token, _ := utils.GenerateToken(user.PublicID.String(), user.Alias)
	refreshToken, _ := utils.GenerateRefreshToken(user.PublicID.String())

	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)
	return utils.Success(ctx, "Login successful", fiber.Map{
		"user":          userResponse,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (c *UserController) GoogleLogin(ctx *fiber.Ctx) error {
	var body struct {
		GoogleID string `json:"google_id"`
		Email    string `json:"email"`
		Alias    string `json:"alias"`
		PhotoURL string `json:"photo_url"`
	}
	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	// Validate required fields
	if body.GoogleID == "" || body.Email == "" {
		return utils.BadRequest(ctx, "Google ID and email are required", "")
	}

	user, err := c.service.GoogleLogin(body.GoogleID, body.Email, body.Alias, body.PhotoURL)
	if err != nil {
		return utils.BadRequest(ctx, "Google login failed", err.Error())
	}

	token, _ := utils.GenerateToken(user.PublicID.String(), user.Alias)
	refreshToken, _ := utils.GenerateRefreshToken(user.PublicID.String())

	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)
	return utils.Success(ctx, "Google login successful", fiber.Map{
		"user":          userResponse,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (c *UserController) UpdateAlias(ctx *fiber.Ctx) error {
	var body struct {
		Alias string `json:"alias"`
	}
	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid request body", err.Error())
	}

	if body.Alias == "" {
		return utils.BadRequest(ctx, "Alias is required", "")
	}

	// Validate alias length
	if len(body.Alias) < 3 || len(body.Alias) > 30 {
		return utils.BadRequest(ctx, "Alias must be between 3 and 30 characters", "")
	}

	// Get public_id from context (set by JWT middleware)
	publicID := ctx.Locals("public_id")
	if publicID == nil {
		return utils.Unauthorized(ctx, "User not authenticated", "")
	}

	// Convert UUID to string
	userIDStr := publicID.(uuid.UUID).String()

	user, err := c.service.UpdateAlias(userIDStr, body.Alias)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to update alias", err.Error())
	}

	var userResponse models.UserResponse
	_ = copier.Copy(&userResponse, &user)
	return utils.Success(ctx, "Alias updated successfully", userResponse)
}
