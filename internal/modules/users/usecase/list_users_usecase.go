package usecase

import (
	"context"

	rulesDto "github.com/dutra/recodationsApp/internal/modules/rules/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/dto"
	"github.com/dutra/recodationsApp/internal/modules/users/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
)

// ListUsersUseCase returns paginated users.
type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: userRepo}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, page, limit int) sharedDomain.Result[[]dto.UserResponse] {
	users, _, err := uc.userRepo.List(ctx, page, limit)
	if err != nil {
		return sharedDomain.Fail[[]dto.UserResponse](err, "failed to list users")
	}

	response := make([]dto.UserResponse, len(users))
	for i, u := range users {
		var rules []rulesDto.RuleResponse
		for _, r := range u.Rules {
			rules = append(rules, rulesDto.RuleResponse{
				ID:          r.ID,
				Name:        r.Name,
				Code:        r.Code,
				Description: r.Description,
				CreatedAt:   r.CreatedAt,
			})
		}

		response[i] = dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Version:   u.Version,
			Rules:     rules,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return sharedDomain.Ok(response, "Users retrieved successfully")
}
