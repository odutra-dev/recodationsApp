package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	auditDomain "github.com/dutra/recodationsApp/internal/modules/audit/domain"
	auditRepo "github.com/dutra/recodationsApp/internal/modules/audit/repository"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/dutra/recodationsApp/pkg/snowflake"
	"github.com/gofiber/fiber/v3"
)

// AuditMiddleware logs modifying HTTP operations (POST, PUT, PATCH, DELETE) to the audit table.
func AuditMiddleware(repo auditRepo.AuditRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		method := c.Method()
		err := c.Next()

		// Only audit mutating actions
		if method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE" {
			path := c.Path()
			// Skip auditing auth login attempts to prevent leaking passwords
			if strings.Contains(path, "/auth/login") {
				return err
			}

			traceID := presenter.GetTraceID(c)
			statusCode := c.Response().StatusCode()
			ip := c.IP()

			var userID *int64
			if uVal := c.Locals("userId"); uVal != nil {
				if uid, ok := uVal.(int64); ok {
					userID = &uid
				}
			}

			resource := path
			action := method

			auditEntry := &auditDomain.AuditLog{
				ID:         snowflake.Generate(),
				UserID:     userID,
				Action:     action,
				Resource:   resource,
				StatusCode: statusCode,
				IPAddress:  ip,
				TraceID:    traceID,
				Details:    string(c.Body()),
				CreatedAt:  time.Now().UTC(),
			}

			go func(entry *auditDomain.AuditLog) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if aErr := repo.Create(ctx, entry); aErr != nil {
					slog.Warn("Failed to persist audit log", "error", aErr)
				}
			}(auditEntry)
		}

		return err
	}
}
