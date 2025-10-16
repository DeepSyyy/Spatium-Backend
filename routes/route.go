package routes

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController, pc *controllers.PostController, cc *controllers.CommentController, rc *controllers.ReactionController) {
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
}
