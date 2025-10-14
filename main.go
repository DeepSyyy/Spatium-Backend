package main

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/DeepSyyy/Spatium-Backend/routes"
	"github.com/DeepSyyy/Spatium-Backend/seed"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load environment variables and connect to the database
	config.LoadEnv()
	config.ConnectDB()

	// Seed initial data
	seed.SeedUsers(config.DB)

	// Initialize Fiber app
	app := fiber.New()

	// Setup for User module
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Setup routes
	routes.Setup(app, userController)

	// Setup port
	port := config.AppConfig.AppPort
	log.Print("Server running on port: " + port)

	// Start the server
	log.Fatal(app.Listen(":" + port))
}
