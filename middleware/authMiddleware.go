package middleware

import (
	"fmt"
	"notes-api/initializers"
	"notes-api/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Validate JWT and Pass user to next function
func AuthHandler(c *fiber.Ctx) error {
	// Get token str from cookie
	tokenString := c.Cookies("Authorization")

	// decode token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Find user from DB
		var user models.User
		// initializers.DB.First(&user, claims["sub"])
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		
		// Attach user to context
		c.Locals("user", user)
		
		return c.Next()
		
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

}