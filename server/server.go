package server

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/DeepSyyy/Spatium-Backend/routes"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/gofiber/fiber/v2"
)

func Start() {
	// Initialize Fiber app
	app := fiber.New()

	// Setup for User module
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// post module setup
	postRepo := repositories.NewPostRepository(config.DB)
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	// comment module setup
	commentRepo := repositories.NewCommentRepository(config.DB)
	commentService := services.NewCommentService(commentRepo)
	commentController := controllers.NewCommentController(commentService)

	// Setup routes
	routes.Setup(app, userController, postController, commentController)

	// Setup port
	port := config.AppConfig.AppPort
	log.Print("Server running on port: " + port)

	// Start the server
	log.Fatal(app.Listen(":" + port))
}
