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

	// ✅ FIX: Always use Railway's PORT when available
	port := os.Getenv("PORT")
	if port == "" {
		// fallback ke APP_PORT untuk local only
		port = config.AppConfig.AppPort
	}

	log.Println("📦 Environment PORT:", os.Getenv("PORT"))
	log.Println("📦 Config APP_PORT:", config.AppConfig.AppPort)
	log.Println("🚀 Server running on port:", port)

	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Printf("❌ Server stopped with error: %v", err)
		log.Println("⚠️ Preventing crash (Railway auto-restart avoided).")
	}
}
