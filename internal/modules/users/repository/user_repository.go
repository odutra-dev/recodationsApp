package repository

import (
	"context"
	"errors"
	"time"

	rulesDomain "github.com/dutra/recodationsApp/internal/modules/rules/domain"
	"github.com/dutra/recodationsApp/internal/modules/users/domain"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"gorm.io/gorm"
)

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, page, limit int) ([]domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
	AssignRules(ctx context.Context, userID int64, ruleIDs []int64) error
}

type postgresUserRepository struct {
	db *gorm.DB
}

// NewUserRepository constructs a GORM UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Rules").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Rules").First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) List(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []domain.User
	err := r.db.WithContext(ctx).
		Preload("Rules").
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, total, err
}

// Update performs an update guarded by Optimistic Locking on the version column.
func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND version = ?", user.ID, user.Version).
		Updates(map[string]any{
			"name":          user.Name,
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"version":       user.Version + 1,
			"updated_at":    now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		var count int64
		_ = r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Count(&count)
		if count == 0 {
			return sharedDomain.ErrNotFound
		}
		return sharedDomain.ErrOptimisticLockConflict
	}

	user.Version++
	user.UpdatedAt = now
	return nil
}

func (r *postgresUserRepository) AssignRules(ctx context.Context, userID int64, ruleIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Clean existing
		if err := tx.Where("user_id = ?", userID).Delete(&rulesDomain.UserRule{}).Error; err != nil {
			return err
		}

		// Insert new associations
		for _, rid := range ruleIDs {
			ur := rulesDomain.UserRule{
				UserID: userID,
				RuleID: rid,
			}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
