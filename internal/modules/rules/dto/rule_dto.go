package dto

import "time"

// CreateRuleInput represents the payload to register a new rule.
type CreateRuleInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// RuleResponse represents the returned rule entity.
type RuleResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
