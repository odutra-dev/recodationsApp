package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/rules/domain"
	"github.com/dutra/recodationsApp/internal/modules/rules/dto"
	"github.com/dutra/recodationsApp/internal/modules/rules/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/pkg/snowflake"
)

// CreateRuleUseCase handles creation of new RBAC rules.
type CreateRuleUseCase struct {
	ruleRepo repository.RuleRepository
}

func NewCreateRuleUseCase(ruleRepo repository.RuleRepository) *CreateRuleUseCase {
	return &CreateRuleUseCase{ruleRepo: ruleRepo}
}

func (uc *CreateRuleUseCase) Execute(ctx context.Context, input dto.CreateRuleInput) sharedDomain.Result[dto.RuleResponse] {
	valErrors := make(map[string][]string)

	name := strings.TrimSpace(input.Name)
	if name == "" {
		valErrors["name"] = append(valErrors["name"], "name is required")
	}

	code := strings.TrimSpace(strings.ToLower(input.Code))
	if code == "" {
		valErrors["code"] = append(valErrors["code"], "code is required")
	}

	if len(valErrors) > 0 {
		return sharedDomain.FailWithValidation[dto.RuleResponse]("validation failed", valErrors)
	}

	// Check if already exists
	existing, err := uc.ruleRepo.FindByCode(ctx, code)
	if err == nil && existing != nil {
		return sharedDomain.Fail[dto.RuleResponse](sharedDomain.ErrAlreadyExists, "rule code already exists")
	}

	rule := &domain.Rule{
		ID:          snowflake.Generate(),
		Name:        name,
		Code:        code,
		Description: input.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := uc.ruleRepo.Create(ctx, rule); err != nil {
		return sharedDomain.Fail[dto.RuleResponse](err, "failed to persist rule")
	}

	response := dto.RuleResponse{
		ID:          rule.ID,
		Name:        rule.Name,
		Code:        rule.Code,
		Description: rule.Description,
		CreatedAt:   rule.CreatedAt,
	}

	return sharedDomain.Ok(response, "Rule created successfully")
}
