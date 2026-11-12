package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
)

// RecoverMiddleware recovers from panics and returns an envelope error.
func RecoverMiddleware() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				traceID := presenter.GetTraceID(c)
				slog.Error("Unhandled panic recovered",
					"error", fmt.Sprintf("%v", r),
					"traceId", traceID,
					"stack", string(debug.Stack()),
				)
				err = presenter.RenderError(c, domain.ErrInternal, "An unexpected server error occurred")
			}
		}()
		return c.Next()
	}
}
