package routes

import (
	"notes-api/handlers"
	"notes-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Authentication/User routes
func SetAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/user")

	authGroup.Post("/signup", limiter.New(middleware.AuthRateConfig), handlers.RegisterHandler)								//Register user
	authGroup.Post("/login", limiter.New(middleware.AuthRateConfig), handlers.LoginHandler)									//User login, returns JWT
	authGroup.Get("/validateToken", limiter.New(middleware.AuthRateConfig), middleware.AuthHandler, handlers.ValidateToken)	//Verify JWT validity

}
