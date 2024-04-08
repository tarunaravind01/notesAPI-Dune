package handlers

import (
	"notes-api/initializers"
	"notes-api/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// validates request body
type requestDTO struct {
	Username    string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


// create new user
// POST
func RegisterHandler(c *fiber.Ctx) error {
	// collect and validate req body
	body := new(requestDTO)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message body",
		})
	}

	// hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost) 
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// save to DB
	user := models.User{Username: body.Username, Password: string(hashedPass)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Created user",
	})

}


// authenticate user and generate JWT
// POST
func LoginHandler(c *fiber.Ctx) error {
	// collect and validate req body
	body := new(requestDTO)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message body",
		})
	}

	//lookup req user
	var user models.User
	initializers.DB.First(&user, "username = ?", body.Username)
	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username/Password incorrect",
		})
	}

	// compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email/Password incorrect",
		})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"uname": user.Username,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	hmacSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error in Token Generation",
		})
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HTTPOnly: true,
		SameSite: "lax",
		Secure: true, //can only be sent over https
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged In Successfully", 
	})
}


// test if JWT is valid
// GET
func ValidateToken(c *fiber.Ctx) error {

	user, ok := c.Locals("user").(models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"UserID": user.ID,
	})

}