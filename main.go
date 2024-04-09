package main

import (
	"fmt"
	"notes-api/initializers"
	"notes-api/middleware"
	"notes-api/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)


func main() {

	// Initialize env and db
	initializers.LoadEnvVariable()
	initializers.ConnectToUserDB()
	initializers.MigrateDB()
	initializers.InitNotesDB()

	// defer closeDB
	defer initializers.CloseDB()

	
	app := fiber.New()

	// fiber middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(helmet.New(middleware.HelmetConfig))

	// all the cookies are encrypted
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_ENC_SECRET"),
	}))
	

	

	// routes
	routes.SetAuthRoutes(app)
	routes.SetNoteRoutes(app)
	

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Unable to collet PORT from env. Defaulting to 3000")
		port = "3000"
	}

	err := app.ListenTLS(":"+port, "./certs/fullChain.pem", "./certs/cert-key.pem")

	if err != nil {
		fmt.Printf("Error starting HTTPS server: %v\n", err)
		fmt.Println("**********Dropping down to HTTP**********")

		// Falling back to http
		if err := app.Listen(":" + port); err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			return
		}
	}
}