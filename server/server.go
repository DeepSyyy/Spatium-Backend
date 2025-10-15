package server

import (
	"log"
	"os"

	"github.com/DeepSyyy/Spatium-Backend/config"
	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/DeepSyyy/Spatium-Backend/routes"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/gofiber/fiber/v2"
)

func Start() {
	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	postRepo := repositories.NewPostRepository(config.DB)
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	commentRepo := repositories.NewCommentRepository(config.DB)
	commentService := services.NewCommentService(commentRepo)
	commentController := controllers.NewCommentController(commentService)

	routes.Setup(app, userController, postController, commentController)

	// ✅ FIX: Prioritize Railway PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = config.AppConfig.AppPort // fallback to .env (local)
	}

	log.Println("🚀 Server running on port:", port)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Fatal(err)
	}
}
