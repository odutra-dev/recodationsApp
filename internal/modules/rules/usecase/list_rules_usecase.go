package usecase

import (
	"context"

	"github.com/dutra/recodationsApp/internal/modules/rules/dto"
	"github.com/dutra/recodationsApp/internal/modules/rules/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
)

// ListRulesUseCase retrieves all registered rules.
type ListRulesUseCase struct {
	ruleRepo repository.RuleRepository
}

func NewListRulesUseCase(ruleRepo repository.RuleRepository) *ListRulesUseCase {
	return &ListRulesUseCase{ruleRepo: ruleRepo}
}

func (uc *ListRulesUseCase) Execute(ctx context.Context) sharedDomain.Result[[]dto.RuleResponse] {
	rules, err := uc.ruleRepo.List(ctx)
	if err != nil {
		return sharedDomain.Fail[[]dto.RuleResponse](err, "failed to fetch rules")
	}

	response := make([]dto.RuleResponse, len(rules))
	for i, r := range rules {
		response[i] = dto.RuleResponse{
			ID:          r.ID,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			CreatedAt:   r.CreatedAt,
		}
	}

	return sharedDomain.Ok(response, "Rules retrieved successfully")
}
