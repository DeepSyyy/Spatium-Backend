package routes

import (
	"log"

	"github.com/DeepSyyy/Spatium-Backend/controllers"
	"github.com/DeepSyyy/Spatium-Backend/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController, pc *controllers.PostController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
}
