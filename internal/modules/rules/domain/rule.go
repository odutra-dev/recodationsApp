package domain

import "time"

// Rule represents a granular permission in the RBAC system.
type Rule struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Code        string    `json:"code" gorm:"size:100;uniqueIndex;not null"`
	Description string    `json:"description" gorm:"size:255"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// UserRule represents the association between a user and a rule.
type UserRule struct {
	UserID int64 `json:"userId" gorm:"primaryKey;column:user_id;index"`
	RuleID int64 `json:"ruleId" gorm:"primaryKey;column:rule_id;index"`
}

func (UserRule) TableName() string {
	return "user_rules"
}
