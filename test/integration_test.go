package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/messages/domain"
	messagesRepo "github.com/dutra/recodationsApp/internal/modules/messages/repository"
	rulesRepo "github.com/dutra/recodationsApp/internal/modules/rules/repository"
	usersRepo "github.com/dutra/recodationsApp/internal/modules/users/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/internal/shared/infrastructure/database"
	"github.com/dutra/recodationsApp/pkg/snowflake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestPostgresIntegrationWithTestcontainers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Spin up real PostgreSQL container using testcontainers-go
	pgContainer, err := tcPostgres.Run(ctx,
		"postgres:16-alpine",
		tcPostgres.WithDatabase("test_recodations"),
		tcPostgres.WithUsername("test_user"),
		tcPostgres.WithPassword("test_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Skipf("Skipping testcontainers integration test (Docker daemon not responsive): %v", err)
		return
	}
	defer func() {
		_ = pgContainer.Terminate(context.Background())
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Run migrations
	err = database.AutoMigrate(db)
	require.NoError(t, err)

	// Run seed
	err = database.SeedInitialData(ctx, db)
	require.NoError(t, err)

	// 1. Test Rules Repository & Seeded Admin
	ruleR := rulesRepo.NewRuleRepository(db)
	userR := usersRepo.NewUserRepository(db)

	rules, err := ruleR.List(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, rules, "Seed should have created default rules")

	admin, err := userR.FindByEmail(ctx, "admin@recodations.com")
	require.NoError(t, err)
	assert.NotNil(t, admin)
	assert.NotEmpty(t, admin.Rules, "Admin should have pre-assigned rules")

	// 2. Test Message Creation & Optimistic Locking
	msgR := messagesRepo.NewMessageRepository(db)
	msg := &domain.Message{
		ID:        snowflake.Generate(),
		Name:      "Integration Test Note",
		Content:   "Testing real Postgres container",
		Status:    domain.MessageStatusDraft,
		Version:   1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = msgR.Create(ctx, msg)
	require.NoError(t, err)

	fetched, err := msgR.FindByID(ctx, msg.ID)
	require.NoError(t, err)
	assert.Equal(t, "Integration Test Note", fetched.Name)
	assert.Equal(t, int64(1), fetched.Version)

	// Perform valid update with matching version
	fetched.Name = "Updated Integration Note"
	err = msgR.Update(ctx, fetched)
	require.NoError(t, err)
	assert.Equal(t, int64(2), fetched.Version)

	// Attempt stale update with old version -> Expect OptimisticLockConflict
	staleMsg := &domain.Message{
		ID:      msg.ID,
		Name:    "Stale Concurrent Update",
		Version: 1, // Stale version! Current version is 2
	}
	err = msgR.Update(ctx, staleMsg)
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedDomain.ErrOptimisticLockConflict)

	// 3. Test State Transition Pattern on real database
	err = fetched.TransitionTo(domain.MessageStatusPublished)
	require.NoError(t, err)
	err = msgR.Update(ctx, fetched)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageStatusPublished, fetched.Status)

	// Invalid transition on real entity: PUBLISHED cannot transition directly to invalid
	err = fetched.TransitionTo(domain.MessageStatus("INVALID_STATUS"))
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedDomain.ErrInvalidStateTransition)

	fmt.Println("All testcontainers integration test scenarios passed successfully!")
}
