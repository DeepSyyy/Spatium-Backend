package routes

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController, pc *controllers.PostController, cc *controllers.CommentController, rc *controllers.ReactionController, chatSessionController *controllers.ChatSessionController, chatMessageController *controllers.ChatMessageController, dailyMoodController *controllers.DailyMoodController, aic *controllers.AIReflectionController, reportController *controllers.ReportController, blockController *controllers.BlockController) {
	if err := godotenv.Load(); err != nil {
		log.Println("🌐 Using environment variables from system (Railway/Production)")
	}

	// =====================================
	// 🔐 AUTH ROUTES
	// =====================================
	app.Post("/api/v1/register", uc.Register)
	app.Post("/api/v1/login", uc.Login)
	app.Post("/api/v1/auth/google", uc.GoogleLogin)
	app.Post("/api/v1/token/refresh", uc.RefreshToken)

	// =====================================
	// 👤 USER PROFILE ROUTES
	// =====================================
	app.Put("/api/v1/user/alias", middlewares.JWTProtected, uc.UpdateAlias)

	// =====================================
	// 📝 POSTS ROUTES
	// =====================================
	app.Post("/api/v1/posts", middlewares.JWTProtected, pc.CreatePost)
	app.Get("/api/v1/posts", middlewares.JWTProtected, pc.GetAllPosts)
	app.Get("/api/v1/posts/:post_id", middlewares.JWTProtected, pc.GetPostDetail)
	app.Get("/api/v1/user/posts", middlewares.JWTProtected, pc.GetPostsByUser)
	app.Put("/api/v1/posts/:post_id", middlewares.JWTProtected, pc.UpdatePost)
	app.Delete("/api/v1/posts/:post_id", middlewares.JWTProtected, pc.DeletePost)

	// =====================================
	// 💬 COMMENTS ROUTES
	// =====================================
	app.Post("/api/v1/comments", middlewares.JWTProtected, cc.CreateComment)
	app.Get("/api/v1/posts/:post_id/comments", middlewares.JWTProtected, cc.GetCommentsByPostID)
	app.Delete("/api/v1/comments/:comment_id", middlewares.JWTProtected, cc.DeleteComment)

	// =====================================
	// ❤️ REACTIONS ROUTES
	// =====================================
	app.Post("/api/v1/posts/:post_id/reactions", middlewares.JWTProtected, rc.ReactToPost)
	app.Get("/api/v1/posts/:post_id/reactions", middlewares.JWTProtected, rc.GetReactionSummary)

	//chat session routes
	// =====================================
	// 💬 SESSION ROUTES
	// =====================================
	app.Post("/api/v1/chat/session", middlewares.JWTProtected, chatSessionController.CreateSession)
	app.Get("/api/v1/chat/session", middlewares.JWTProtected, chatSessionController.GetUserSessions)
	app.Delete("/api/v1/chat/session/:id", middlewares.JWTProtected, chatSessionController.DeleteSession)

	// =====================================
	// 💭 MESSAGE ROUTES
	// =====================================
	app.Post("/api/v1/chat/:session_id", middlewares.JWTProtected, chatMessageController.Create)
	app.Get("/api/v1/chat/:session_id/messages", middlewares.JWTProtected, chatMessageController.GetMessagesBySession)
	app.Get("/api/v1/chat/:session_id/last", middlewares.JWTProtected, chatMessageController.GetLastMessages)

	// =====================================
	// 💭 DAILY MOODS ROUTES
	// =====================================
	app.Post("/api/v1/moods", middlewares.JWTProtected, dailyMoodController.CreateOrUpdate)
	app.Get("/api/v1/moods/today", middlewares.JWTProtected, dailyMoodController.GetTodayMood)
	app.Get("/api/v1/moods/weekly", middlewares.JWTProtected, dailyMoodController.GetWeeklyMoods)
	app.Get("/api/v1/moods/statistics", middlewares.JWTProtected, dailyMoodController.GetMoodStatistics)
	app.Get("/api/v1/moods/chart", middlewares.JWTProtected, dailyMoodController.GetMoodChart)

	// =====================================
	// 🧠 AI REFLECTION
	// =====================================
	app.Post("/api/v1/ai/reflection", middlewares.JWTProtected, aic.GenerateReflection)
	app.Get("/api/v1/ai/reflection", middlewares.JWTProtected, aic.GetUserReflection)

	// =====================================
	// 🚨 REPORTS ROUTES (Content Moderation)
	// =====================================
	app.Get("/api/v1/reports/reasons", middlewares.JWTProtected, reportController.GetReportReasons)
	app.Post("/api/v1/reports", middlewares.JWTProtected, reportController.CreateReport)
	app.Get("/api/v1/reports/me", middlewares.JWTProtected, reportController.GetMyReports)
	// Admin routes for report management
	app.Get("/api/v1/admin/reports/pending", middlewares.JWTProtected, reportController.GetPendingReports)
	app.Put("/api/v1/admin/reports/:report_id/status", middlewares.JWTProtected, reportController.UpdateReportStatus)

	// =====================================
	// 🚫 BLOCK ROUTES (User Safety)
	// =====================================
	app.Post("/api/v1/users/block", middlewares.JWTProtected, blockController.BlockUser)
	app.Delete("/api/v1/users/block/:user_id", middlewares.JWTProtected, blockController.UnblockUser)
	app.Get("/api/v1/users/blocked", middlewares.JWTProtected, blockController.GetBlockedUsers)
}
