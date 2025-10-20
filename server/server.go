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

	// Repositories
	userRepo := repositories.NewUserRepository()
	commentRepo := repositories.NewCommentRepository(config.DB)
	chatSessionRepo := repositories.NewChatSessionRepository()
	chatMessageRepo := repositories.NewChatMessageRepository()
	postRepo := repositories.NewPostRepository(config.DB)
	reactionRepo := repositories.NewReactionRepository(config.DB)
	dailyMoodRepo := repositories.NewDailyMoodRepository(config.DB)
	aireflectionRepo := repositories.NewAIReflectionRepository(config.DB)

	// Services
	userService := services.NewUserService(userRepo)
	commentService := services.NewCommentService(commentRepo)
	chatSessionService := services.NewChatSessionService(chatSessionRepo)
	chatMessageService := services.NewChatMessageService(chatMessageRepo, chatSessionRepo)
	postService := services.NewPostService(postRepo)
	reactionService := services.NewReactionService(reactionRepo)
	aiReflectionService := services.NewAIReflectionService(aireflectionRepo)
	dailyMoodService := services.NewDailyMoodService(dailyMoodRepo, aiReflectionService)

	// Controllers
	userController := controllers.NewUserController(userService)
	chatSessionController := controllers.NewChatSessionController(chatSessionService)
	chatMessageController := controllers.NewChatMessageController(chatMessageService)
	commentController := controllers.NewCommentController(commentService, postService)
	postController := controllers.NewPostController(postService, commentService)
	reactionController := controllers.NewReactionController(reactionService, postService)
	dailyMoodController := controllers.NewDailyMoodController(dailyMoodService)
	aiReflectionController := controllers.NewAIReflectionController(aiReflectionService)

	routes.Setup(app, userController, postController, commentController, reactionController, chatSessionController, chatMessageController, dailyMoodController, aiReflectionController)

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
