package usecase

import (
	"context"
	"fmt"

	"github.com/dutra/recodationsApp/internal/modules/users/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/cache"
)

// AssignRulesUseCase updates a user's associated permission rules and clears cache.
type AssignRulesUseCase struct {
	userRepo repository.UserRepository
	cache    cache.Cache
}

func NewAssignRulesUseCase(userRepo repository.UserRepository, cache cache.Cache) *AssignRulesUseCase {
	return &AssignRulesUseCase{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (uc *AssignRulesUseCase) Execute(ctx context.Context, userID int64, input dto.AssignRulesInput) sharedDomain.Result[string] {
	if userID <= 0 {
		return sharedDomain.Fail[string](sharedDomain.ErrValidation, "invalid user id")
	}

	// Verify user exists
	_, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return sharedDomain.Fail[string](err, "user not found")
	}

	if err := uc.userRepo.AssignRules(ctx, userID, input.RuleIDs); err != nil {
		return sharedDomain.Fail[string](err, "failed to assign rules")
	}

	// Invalidate L1/L2 Cache
	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, fmt.Sprintf("user:%d", userID))
	}

	return sharedDomain.Ok("Rules assigned successfully", "Success")
}
