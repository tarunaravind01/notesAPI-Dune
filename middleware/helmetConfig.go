package middleware

import "github.com/gofiber/fiber/v2/middleware/helmet"


// Adding Security headers using helmet
var HelmetConfig = helmet.Config{
    HSTSExcludeSubdomains: false,
    HSTSPreloadEnabled:    true,
    CrossOriginResourcePolicy: "same-origin",
}
