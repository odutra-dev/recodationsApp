package middleware

import (
	"time"

	"github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// RateLimitMiddleware provides protection against traffic spikes and brute force attacks.
func RateLimitMiddleware(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c fiber.Ctx) error {
			return presenter.RenderResponse(c, fiber.StatusTooManyRequests, domain.Fail[any](
				domain.ErrForbidden,
				"Too many requests. Please wait before retrying.",
			))
		},
	})
}
