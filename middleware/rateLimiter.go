package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

//rate limiter config for Auth functions
//the rate can be reduced or increased
var AuthRateConfig = limiter.Config{
	Max: 		20,
	Expiration: 30 * time.Second,
	LimiterMiddleware: limiter.SlidingWindow{},
	LimitReached: rateLimitExceeded,

}


//rate limiter config for notes endpoint
var NotesRateConfig = limiter.Config{
	Max: 		45,
	Expiration: 30 * time.Second,
	LimiterMiddleware: limiter.SlidingWindow{},
	LimitReached: rateLimitExceeded,

}

func rateLimitExceeded(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Slow Down! Rate Limit Reached",
	})
}

