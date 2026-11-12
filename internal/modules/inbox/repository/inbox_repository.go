package repository

import (
	"context"
	"errors"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/inbox/domain"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"gorm.io/gorm"
)

// InboxRepository defines persistence operations for idempotent inbox pattern.
type InboxRepository interface {
	Create(ctx context.Context, msg *domain.InboxMessage) error
	FindByEventID(ctx context.Context, eventID string) (*domain.InboxMessage, error)
	FindPending(ctx context.Context, limit int) ([]domain.InboxMessage, error)
	Update(ctx context.Context, msg *domain.InboxMessage) error
}

type postgresInboxRepository struct {
	db *gorm.DB
}

// NewInboxRepository constructs a GORM InboxRepository.
func NewInboxRepository(db *gorm.DB) InboxRepository {
	return &postgresInboxRepository{db: db}
}

func (r *postgresInboxRepository) Create(ctx context.Context, msg *domain.InboxMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *postgresInboxRepository) FindByEventID(ctx context.Context, eventID string) (*domain.InboxMessage, error) {
	var msg domain.InboxMessage
	err := r.db.WithContext(ctx).First(&msg, "event_id = ?", eventID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return &msg, nil
}

func (r *postgresInboxRepository) FindPending(ctx context.Context, limit int) ([]domain.InboxMessage, error) {
	var messages []domain.InboxMessage
	err := r.db.WithContext(ctx).
		Where("status = ?", domain.InboxStatusPending).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// Update guards state changes with Optimistic Locking.
func (r *postgresInboxRepository) Update(ctx context.Context, msg *domain.InboxMessage) error {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Model(&domain.InboxMessage{}).
		Where("id = ? AND version = ?", msg.ID, msg.Version).
		Updates(map[string]any{
			"status":        msg.Status,
			"error_message": msg.ErrorMessage,
			"attempts":      msg.Attempts,
			"processed_at":  msg.ProcessedAt,
			"version":       msg.Version + 1,
			"updated_at":    now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		var count int64
		_ = r.db.WithContext(ctx).Model(&domain.InboxMessage{}).Where("id = ?", msg.ID).Count(&count)
		if count == 0 {
			return sharedDomain.ErrNotFound
		}
		return sharedDomain.ErrOptimisticLockConflict
	}

	msg.Version++
	msg.UpdatedAt = now
	return nil
}
