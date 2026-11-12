package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	auditDomain "github.com/dutra/recodationsApp/internal/modules/audit/domain"
	inboxDomain "github.com/dutra/recodationsApp/internal/modules/inbox/domain"
	messagesDomain "github.com/dutra/recodationsApp/internal/modules/messages/domain"
	rulesDomain "github.com/dutra/recodationsApp/internal/modules/rules/domain"
	usersDomain "github.com/dutra/recodationsApp/internal/modules/users/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds PostgreSQL connection parameters.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connect establishes a connection to PostgreSQL and runs schema automigrations.
func Connect(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	gormLogger := logger.Default.LogMode(logger.Silent)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sql.DB: %w", err)
	}

	// Connection pool resilience settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run automigrations
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("automigration failed: %w", err)
	}

	// Run initial seed
	if err := SeedInitialData(context.Background(), db); err != nil {
		slog.Warn("Failed running initial seed", "error", err)
	}

	return db, nil
}

// AutoMigrate migrates all domain schemas.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&rulesDomain.Rule{},
		&rulesDomain.UserRule{},
		&usersDomain.User{},
		&messagesDomain.Message{},
		&inboxDomain.InboxMessage{},
		&auditDomain.AuditLog{},
	)
}
