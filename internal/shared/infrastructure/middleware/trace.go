package middleware

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"
)

// TraceMiddleware ensures every HTTP request has a unique TraceID for observability.
func TraceMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Locals("traceId", traceID)
		c.Set("X-Trace-ID", traceID)

		return c.Next()
	}
}
