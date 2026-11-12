package domain

import (
	"time"

	rulesDomain "github.com/dutra/recodationsApp/internal/modules/rules/domain"
)

// User represents a system user with RBAC rules and optimistic concurrency lock.
type User struct {
	ID           int64              `json:"id" gorm:"primaryKey"`
	Name         string             `json:"name" gorm:"size:255;not null"`
	Email        string             `json:"email" gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string             `json:"-" gorm:"size:255;not null"`
	Version      int64              `json:"version" gorm:"default:1;not null"`
	Rules        []rulesDomain.Rule `json:"rules,omitempty" gorm:"many2many:user_rules;"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
