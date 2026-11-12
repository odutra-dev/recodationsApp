package repository

import (
	"context"
	"errors"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/messages/domain"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"gorm.io/gorm"
)

// MessageRepository defines persistence operations for messages.
type MessageRepository interface {
	Create(ctx context.Context, msg *domain.Message) error
	FindByID(ctx context.Context, id int64) (*domain.Message, error)
	List(ctx context.Context, page, limit int) ([]domain.Message, int64, error)
	Update(ctx context.Context, msg *domain.Message) error
	Delete(ctx context.Context, id int64) error
}

type postgresMessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository constructs a GORM MessageRepository.
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &postgresMessageRepository{db: db}
}

func (r *postgresMessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *postgresMessageRepository) FindByID(ctx context.Context, id int64) (*domain.Message, error) {
	var msg domain.Message
	err := r.db.WithContext(ctx).First(&msg, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &msg, nil
}

func (r *postgresMessageRepository) List(ctx context.Context, page, limit int) ([]domain.Message, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.Message{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var messages []domain.Message
	err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error

	return messages, total, err
}

// Update updates the message with Optimistic Locking verification.
func (r *postgresMessageRepository) Update(ctx context.Context, msg *domain.Message) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&domain.Message{}).
		Where("id = ? AND version = ?", msg.ID, msg.Version).
		Updates(map[string]any{
			"name":       msg.Name,
			"content":    msg.Content,
			"image":      msg.Image,
			"status":     msg.Status,
			"version":    msg.Version + 1,
			"updated_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		var count int64
		_ = r.db.WithContext(ctx).Model(&domain.Message{}).Where("id = ?", msg.ID).Count(&count)
		if count == 0 {
			return sharedDomain.ErrNotFound
		}
		return sharedDomain.ErrOptimisticLockConflict
	}

	msg.Version++
	msg.UpdatedAt = now
	return nil
}

func (r *postgresMessageRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&domain.Message{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sharedDomain.ErrNotFound
	}
	return nil
}
