package presenter

import (
	"errors"
	"time"

	"github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/gofiber/fiber/v3"
)

// ResponseHTTP matches the contract requested by the user:
//
//	public record ResponseHttp<T>(
//	    T? Data,
//	    string Message,
//	    string TraceId,
//	    bool Success,
//	    DateTimeOffset Timestamp,
//	    string? DetailsError = null,
//	    string? Path = null,
//	    UInt16 ApiVersion = 1,
//	    IReadOnlyList<string>? Errors = null,
//	    IDictionary<string, string[]>? ErrorValids = null
//	);
type ResponseHTTP[T any] struct {
	Data         *T                  `json:"data"`
	Message      string              `json:"message"`
	TraceID      string              `json:"traceId"`
	Success      bool                `json:"success"`
	Timestamp    time.Time           `json:"timestamp"`
	DetailsError *string             `json:"detailsError,omitempty"`
	Path         *string             `json:"path,omitempty"`
	APIVersion   uint16              `json:"apiVersion"`
	Errors       []string            `json:"errors,omitempty"`
	ErrorValids  map[string][]string `json:"errorValids,omitempty"`
}

// GetTraceID extracts the trace ID from Fiber Context locals, falling back to empty string.
func GetTraceID(c fiber.Ctx) string {
	if val, ok := c.Locals("traceId").(string); ok && val != "" {
		return val
	}
	return ""
}

// RenderResponse builds and sends the standard ResponseHTTP[T] envelope.
func RenderResponse[T any](c fiber.Ctx, statusCode int, result domain.Result[T]) error {
	traceID := GetTraceID(c)
	reqPath := c.Path()

	var details *string
	if result.DetailsError != "" {
		details = &result.DetailsError
	}

	response := ResponseHTTP[T]{
		Data:         result.Data,
		Message:      result.Message,
		TraceID:      traceID,
		Success:      result.Success,
		Timestamp:    time.Now().UTC(),
		DetailsError: details,
		Path:         &reqPath,
		APIVersion:   1,
		Errors:       result.Errors,
		ErrorValids:  result.ErrorValids,
	}

	return c.Status(statusCode).JSON(response)
}

// RenderSuccess returns a successful HTTP response with data.
func RenderSuccess[T any](c fiber.Ctx, statusCode int, data T, message string) error {
	res := domain.Ok(data, message)
	return RenderResponse(c, statusCode, res)
}

// RenderEmpty returns a successful HTTP response without data.
func RenderEmpty(c fiber.Ctx, statusCode int, message string) error {
	res := domain.OkEmpty[any](message)
	return RenderResponse(c, statusCode, res)
}

// RenderError maps domain errors to appropriate HTTP status codes and responds with ResponseHTTP[any].
func RenderError(c fiber.Ctx, err error, message string) error {
	statusCode := MapErrorToStatusCode(err)
	res := domain.Fail[any](err, message)
	return RenderResponse(c, statusCode, res)
}

// RenderValidationError returns an HTTP 400 Bad Request with validation breakdown.
func RenderValidationError(c fiber.Ctx, message string, valErrors map[string][]string) error {
	res := domain.FailWithValidation[any](message, valErrors)
	return RenderResponse(c, fiber.StatusBadRequest, res)
}

// MapErrorToStatusCode maps domain error types to HTTP status codes.
func MapErrorToStatusCode(err error) int {
	if err == nil {
		return fiber.StatusOK
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, domain.ErrAlreadyExists):
		return fiber.StatusConflict
	case errors.Is(err, domain.ErrOptimisticLockConflict):
		return fiber.StatusConflict
	case errors.Is(err, domain.ErrIdempotencyConflict):
		return fiber.StatusConflict
	case errors.Is(err, domain.ErrValidation):
		return fiber.StatusBadRequest
	case errors.Is(err, domain.ErrInvalidStateTransition):
		return fiber.StatusUnprocessableEntity
	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidCredentials):
		return fiber.StatusUnauthorized
	case errors.Is(err, domain.ErrForbidden):
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}
