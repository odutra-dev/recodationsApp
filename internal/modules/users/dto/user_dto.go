package dto

import (
	"time"

	rulesDto "github.com/dutra/recodationsApp/internal/modules/rules/dto"
)

// CreateUserInput represents data required to create a user.
type CreateUserInput struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	RuleIDs  []int64 `json:"ruleIds,omitempty"`
}

// AssignRulesInput allows updating a user's permission rules.
type AssignRulesInput struct {
	RuleIDs []int64 `json:"ruleIds"`
}

// UserResponse represents the public user output.
type UserResponse struct {
	ID        int64                  `json:"id"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Version   int64                  `json:"version"`
	Rules     []rulesDto.RuleResponse `json:"rules,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}
