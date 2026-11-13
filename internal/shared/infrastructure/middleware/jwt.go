package middleware

import (
	"errors"
	"strings"

	"github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/presenter"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID int64    `json:"user_id"`
	Email  string   `json:"email"`
	Rules  []string `json:"rules"`
	jwt.RegisteredClaims
}

// JWTMiddleware authenticates requests by validating Bearer JWT tokens.
func JWTMiddleware(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return presenter.RenderError(c, domain.ErrUnauthorized, "Missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return presenter.RenderError(c, domain.ErrUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'")
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return presenter.RenderError(c, domain.ErrUnauthorized, "Invalid or expired token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return presenter.RenderError(c, domain.ErrUnauthorized, "Malformed token claims")
		}

		c.Locals("userId", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("rules", claims.Rules)

		return c.Next()
	}
}

// RequireRule enforces that the authenticated user possesses a specific rule/permission.
func RequireRule(requiredRule string) fiber.Handler {
	return func(c fiber.Ctx) error {
		rulesRaw := c.Locals("rules")
		if rulesRaw == nil {
			return presenter.RenderError(c, domain.ErrForbidden, "No permissions associated with this user")
		}

		rules, ok := rulesRaw.([]string)
		if !ok {
			return presenter.RenderError(c, domain.ErrForbidden, "Invalid permissions payload")
		}

		for _, r := range rules {
			if r == requiredRule || r == "admin" || r == "*" {
				return c.Next()
			}
		}

		return presenter.RenderError(c, domain.ErrForbidden, "You do not have the required permission: "+requiredRule)
	}
}
