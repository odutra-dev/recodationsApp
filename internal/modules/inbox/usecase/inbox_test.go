package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/inbox/domain"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockInboxRepo struct {
	mock.Mock
}

func (m *mockInboxRepo) Create(ctx context.Context, msg *domain.InboxMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockInboxRepo) FindByEventID(ctx context.Context, eventID string) (*domain.InboxMessage, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InboxMessage), args.Error(1)
}

func (m *mockInboxRepo) FindPending(ctx context.Context, limit int) ([]domain.InboxMessage, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]domain.InboxMessage), args.Error(1)
}

func (m *mockInboxRepo) Update(ctx context.Context, msg *domain.InboxMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func TestInboxSaveEventSuccess(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := NewInboxService(repo)

	ctx := context.Background()
	eventID := "evt-12345"

	repo.On("FindByEventID", ctx, eventID).Return(nil, sharedDomain.ErrNotFound)
	repo.On("Create", ctx, mock.AnythingOfType("*domain.InboxMessage")).Return(nil)

	err := svc.SaveEvent(ctx, eventID, "message.created", map[string]string{"foo": "bar"})
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestInboxSaveEventDuplicateIdempotency(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := NewInboxService(repo)

	ctx := context.Background()
	eventID := "evt-existing"

	existing := &domain.InboxMessage{
		ID:        1,
		EventID:   eventID,
		Status:    domain.InboxStatusProcessed,
		CreatedAt: time.Now(),
	}

	repo.On("FindByEventID", ctx, eventID).Return(existing, nil)

	err := svc.SaveEvent(ctx, eventID, "message.created", map[string]string{"foo": "bar"})
	require.Error(t, err)
	assert.ErrorIs(t, err, sharedDomain.ErrIdempotencyConflict)
}

func TestInboxProcessPendingEvents(t *testing.T) {
	repo := new(mockInboxRepo)
	svc := NewInboxService(repo)

	ctx := context.Background()
	pendingList := []domain.InboxMessage{
		{
			ID:        1,
			EventID:   "evt-1",
			EventType: "order.created",
			Status:    domain.InboxStatusPending,
			Version:   1,
		},
		{
			ID:        2,
			EventID:   "evt-2",
			EventType: "order.failed",
			Status:    domain.InboxStatusPending,
			Version:   1,
		},
	}

	repo.On("FindPending", ctx, 10).Return(pendingList, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(m *domain.InboxMessage) bool {
		return m.EventID == "evt-1" && m.Status == domain.InboxStatusProcessed
	})).Return(nil)
	repo.On("Update", ctx, mock.MatchedBy(func(m *domain.InboxMessage) bool {
		return m.EventID == "evt-2" && m.Status == domain.InboxStatusFailed
	})).Return(nil)

	count, err := svc.ProcessPendingEvents(ctx, 10, func(c context.Context, event domain.InboxMessage) error {
		if event.EventID == "evt-2" {
			return errors.New("simulated handler error")
		}
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, count)
	repo.AssertExpectations(t)
}
