package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{service: userService}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user, err := c.service.Register()
	if err != nil {
		// Handle error (e.g., return HTTP 500)
		return utils.BadRequest(nil, "Registration failed", err.Error())
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
