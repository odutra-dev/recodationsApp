package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/auth/dto"
	rulesDto "github.com/dutra/recodationsApp/internal/modules/rules/dto"
	usersDto "github.com/dutra/recodationsApp/internal/modules/users/dto"
	usersRepo "github.com/dutra/recodationsApp/internal/modules/users/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/middleware"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginUseCase authenticates users and issues signed JWT tokens.
type LoginUseCase struct {
	userRepo  usersRepo.UserRepository
	jwtSecret string
	ttl       time.Duration
}

func NewLoginUseCase(userRepo usersRepo.UserRepository, jwtSecret string, ttl time.Duration) *LoginUseCase {
	return &LoginUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		ttl:       ttl,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input dto.LoginInput) sharedDomain.Result[dto.LoginResponse] {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return sharedDomain.Fail[dto.LoginResponse](sharedDomain.ErrValidation, "email and password are required")
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return sharedDomain.Fail[dto.LoginResponse](sharedDomain.ErrInvalidCredentials, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return sharedDomain.Fail[dto.LoginResponse](sharedDomain.ErrInvalidCredentials, "invalid email or password")
	}

	// Collect rules
	var ruleCodes []string
	var rulesResponse []rulesDto.RuleResponse
	for _, r := range user.Rules {
		ruleCodes = append(ruleCodes, r.Code)
		rulesResponse = append(rulesResponse, rulesDto.RuleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
		})
	}

	now := time.Now().UTC()
	expiresAt := now.Add(uc.ttl)

	claims := middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Rules:  ruleCodes,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return sharedDomain.Fail[dto.LoginResponse](err, "failed generating access token")
	}

	response := dto.LoginResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(uc.ttl.Seconds()),
		User: usersDto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Version:   user.Version,
			Rules:     rulesResponse,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return sharedDomain.Ok(response, "Login successful")
}
