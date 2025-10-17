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

	// Initialize repositories, services, and controllers
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Initialize Post and Comment components
	commentRepo := repositories.NewCommentRepository(config.DB)
	commentService := services.NewCommentService(commentRepo)
	postRepo := repositories.NewPostRepository(config.DB)
	postService := services.NewPostService(postRepo)

	// Initialize Comment and post controller
	postController := controllers.NewPostController(postService, commentService)
	commentController := controllers.NewCommentController(commentService, postService)

	// initialize reaction controller
	reactionRepo := repositories.NewReactionRepository(config.DB)
	reactionService := services.NewReactionService(reactionRepo)
	reactionController := controllers.NewReactionController(reactionService, postService)

	// Repositories
	chatSessionRepo := repositories.NewChatSessionRepository()
	chatMessageRepo := repositories.NewChatMessageRepository()

	// Services
	chatSessionService := services.NewChatSessionService(chatSessionRepo)
	chatMessageService := services.NewChatMessageService(chatMessageRepo, chatSessionRepo)

	// Controllers
	chatSessionController := controllers.NewChatSessionController(chatSessionService)
	chatMessageController := controllers.NewChatMessageController(chatMessageService)

	routes.Setup(app, userController, postController, commentController, reactionController, chatSessionController, chatMessageController)

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
