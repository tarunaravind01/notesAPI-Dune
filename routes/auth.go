package routes

import (
	"notes-api/handlers"
	"notes-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// Authentication/User routes
func SetAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/user")

	authGroup.Post("/signup", handlers.RegisterHandler)								//Register user
	authGroup.Post("/login", handlers.LoginHandler)									//User login, returns JWT
	authGroup.Get("/validateToken", middleware.AuthHandler, handlers.ValidateToken)	//Verify JWT validity

}
