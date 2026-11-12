package repository

import (
	"context"
	"errors"

	"github.com/dutra/recodationsApp/internal/modules/rules/domain"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"gorm.io/gorm"
)

// RuleRepository defines persistence operations for permissions.
type RuleRepository interface {
	Create(ctx context.Context, rule *domain.Rule) error
	FindByID(ctx context.Context, id int64) (*domain.Rule, error)
	FindByCode(ctx context.Context, code string) (*domain.Rule, error)
	List(ctx context.Context) ([]domain.Rule, error)
	FindRulesByUserID(ctx context.Context, userID int64) ([]domain.Rule, error)
	AssignRuleToUser(ctx context.Context, userID int64, ruleID int64) error
	RemoveRuleFromUser(ctx context.Context, userID int64, ruleID int64) error
}

type postgresRuleRepository struct {
	db *gorm.DB
}

// NewRuleRepository constructs a GORM-backed RuleRepository.
func NewRuleRepository(db *gorm.DB) RuleRepository {
	return &postgresRuleRepository{db: db}
}

func (r *postgresRuleRepository) Create(ctx context.Context, rule *domain.Rule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

func (r *postgresRuleRepository) FindByID(ctx context.Context, id int64) (*domain.Rule, error) {
	var rule domain.Rule
	err := r.db.WithContext(ctx).First(&rule, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *postgresRuleRepository) FindByCode(ctx context.Context, code string) (*domain.Rule, error) {
	var rule domain.Rule
	err := r.db.WithContext(ctx).First(&rule, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *postgresRuleRepository) List(ctx context.Context) ([]domain.Rule, error) {
	var rules []domain.Rule
	err := r.db.WithContext(ctx).Order("code ASC").Find(&rules).Error
	return rules, err
}

func (r *postgresRuleRepository) FindRulesByUserID(ctx context.Context, userID int64) ([]domain.Rule, error) {
	var rules []domain.Rule
	err := r.db.WithContext(ctx).
		Table("rules").
		Joins("JOIN user_rules ON user_rules.rule_id = rules.id").
		Where("user_rules.user_id = ?", userID).
		Find(&rules).Error
	return rules, err
}

func (r *postgresRuleRepository) AssignRuleToUser(ctx context.Context, userID int64, ruleID int64) error {
	userRule := domain.UserRule{
		UserID: userID,
		RuleID: ruleID,
	}
	return r.db.WithContext(ctx).FirstOrCreate(&userRule).Error
}

func (r *postgresRuleRepository) RemoveRuleFromUser(ctx context.Context, userID int64, ruleID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND rule_id = ?", userID, ruleID).
		Delete(&domain.UserRule{}).Error
}
