package usecase

import (
	"context"
	"fmt"
	"time"

	rulesDto "github.com/dutra/recodationsApp/internal/modules/rules/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
)

// GetUserUseCase retrieves user details, leveraging the Two-Level Cache.
type GetUserUseCase struct {
	userRepo repository.UserRepository
	cache    cache.Cache
}

func NewGetUserUseCase(userRepo repository.UserRepository, cache cache.Cache) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id int64) sharedDomain.Result[dto.UserResponse] {
	if id <= 0 {
		return sharedDomain.Fail[dto.UserResponse](sharedDomain.ErrValidation, "invalid user id")
	}

	cacheKey := fmt.Sprintf("user:%d", id)

	// 1. Check L1/L2 Cache
	if uc.cache != nil {
		var cached dto.UserResponse
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil && cached.ID != 0 {
			return sharedDomain.Ok(cached, "User retrieved from cache")
		}
	}

	// 2. Query Repository
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return sharedDomain.Fail[dto.UserResponse](err, "user not found")
	}

	var rules []rulesDto.RuleResponse
	for _, r := range user.Rules {
		rules = append(rules, rulesDto.RuleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
		})
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Version:   user.Version,
		Rules:     rules,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 3. Populate L1/L2 Cache with TTL
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, cacheKey, response, 10*time.Minute)
	}

	return sharedDomain.Ok(response, "User retrieved successfully")
}
