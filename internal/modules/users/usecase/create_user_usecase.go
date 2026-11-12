package usecase

import (
	"context"
	"strings"
	"time"

	rulesDto "github.com/dutra/recodationsApp/internal/modules/rules/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/domain"
	"github.com/dutra/recodationsApp/internal/modules/users/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserUseCase handles registration of a user.
type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) sharedDomain.Result[dto.UserResponse] {
	valErrors := make(map[string][]string)

	name := strings.TrimSpace(input.Name)
	if name == "" {
		valErrors["name"] = append(valErrors["name"], "name is required")
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || !strings.Contains(email, "@") {
		valErrors["email"] = append(valErrors["email"], "valid email is required")
	}

	if len(input.Password) < 6 {
		valErrors["password"] = append(valErrors["password"], "password must be at least 6 characters long")
	}

	if len(valErrors) > 0 {
		return sharedDomain.FailWithValidation[dto.UserResponse]("validation failed", valErrors)
	}

	// Uniqueness check
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return sharedDomain.Fail[dto.UserResponse](sharedDomain.ErrAlreadyExists, "email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return sharedDomain.Fail[dto.UserResponse](err, "failed to encrypt password")
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           snowflake.Generate(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Version:      1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return sharedDomain.Fail[dto.UserResponse](err, "failed to persist user")
	}

	if len(input.RuleIDs) > 0 {
		_ = uc.userRepo.AssignRules(ctx, user.ID, input.RuleIDs)
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Version:   user.Version,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Rules:     []rulesDto.RuleResponse{},
	}

	return sharedDomain.Ok(response, "User registered successfully")
}
