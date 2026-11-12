package repository

import (
	"context"

	"github.com/dutra/recodationsApp/internal/modules/audit/domain"
	"gorm.io/gorm"
)

// AuditRepository defines data operations for system audit trails.
type AuditRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	List(ctx context.Context, page, limit int) ([]domain.AuditLog, int64, error)
}

type postgresAuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository constructs the GORM implementation of AuditRepository.
func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &postgresAuditRepository{db: db}
}

func (r *postgresAuditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *postgresAuditRepository) List(ctx context.Context, page, limit int) ([]domain.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []domain.AuditLog
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, total, err
}
