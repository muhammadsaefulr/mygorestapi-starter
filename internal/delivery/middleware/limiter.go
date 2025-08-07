package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/muhammadsaefulr/mygorestapi-starter/internal/domain/dto/util/response"
)

func LimiterConfig() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,
		Expiration: 15 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(response.Common{
					Code:    fiber.StatusTooManyRequests,
					Status:  "error",
					Message: "Too many requests, please try again later",
				})
		},
		SkipSuccessfulRequests: true,
	})
}
