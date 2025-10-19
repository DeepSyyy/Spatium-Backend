package routes

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController, pc *controllers.PostController, cc *controllers.CommentController, rc *controllers.ReactionController, chatSessionController *controllers.ChatSessionController, chatMessageController *controllers.ChatMessageController, dailyMoodController *controllers.DailyMoodController) {
	if err := godotenv.Load(); err != nil {
		log.Println("🌐 Using environment variables from system (Railway/Production)")
	}

	//auth routes
	app.Post("/api/v1/register", uc.Register)
	app.Post("/api/v1/login", uc.Login)

	//post routes
	app.Post("/api/v1/posts", middlewares.JWTProtected, pc.CreatePost)
	app.Get("/api/v1/posts", middlewares.JWTProtected, pc.GetAllPosts)
	app.Get("/api/v1/posts/:post_id", middlewares.JWTProtected, pc.GetPostDetail)
	app.Get("/api/v1/user/posts", middlewares.JWTProtected, pc.GetPostsByUser)
	app.Put("/api/v1/posts/:post_id", middlewares.JWTProtected, pc.UpdatePost)
	app.Delete("/api/v1/posts/:post_id", middlewares.JWTProtected, pc.DeletePost)

	//comment routes
	app.Post("/api/v1/comments", middlewares.JWTProtected, cc.CreateComment)
	app.Get("/api/v1/posts/:post_id/comments", middlewares.JWTProtected, cc.GetCommentsByPostID)
	app.Delete("/api/v1/comments/:comment_id", middlewares.JWTProtected, cc.DeleteComment)

	//reaction routes
	app.Post("/api/v1/posts/:post_id/reactions", middlewares.JWTProtected, rc.ReactToPost)
	app.Get("/api/v1/posts/:post_id/reactions", middlewares.JWTProtected, rc.GetReactionSummary)

	//chat session routes
	// =====================================
	// 💬 SESSION ROUTES
	// =====================================
	app.Post("/api/v1/chat/session", middlewares.JWTProtected, chatSessionController.CreateSession)       // create session
	app.Get("/api/v1/chat/session", middlewares.JWTProtected, chatSessionController.GetUserSessions)      // get all sessions
	app.Delete("/api/v1/chat/session/:id", middlewares.JWTProtected, chatSessionController.DeleteSession) // delete session

	// =====================================
	// 💭 MESSAGE ROUTES
	// =====================================
	app.Post("/api/v1/chat/:session_id", middlewares.JWTProtected, chatMessageController.Create)                       // send message (user)
	app.Get("/api/v1/chat/:session_id/messages", middlewares.JWTProtected, chatMessageController.GetMessagesBySession) // get all messages
	app.Get("/api/v1/chat/:session_id/last", middlewares.JWTProtected, chatMessageController.GetLastMessages)          // get last messages

	//daily mood routes
	app.Post("/api/v1/moods", middlewares.JWTProtected, dailyMoodController.CreateOrUpdate)
	app.Get("/api/v1/moods/today", middlewares.JWTProtected, dailyMoodController.GetTodayMood)
	app.Get("/api/v1/moods/weekly", middlewares.JWTProtected, dailyMoodController.GetWeeklyMoods)

}
