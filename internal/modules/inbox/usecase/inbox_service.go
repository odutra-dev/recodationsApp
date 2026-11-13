package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/dutra/recodationsApp/internal/modules/inbox/domain"
	"github.com/dutra/recodationsApp/internal/modules/inbox/repository"
	sharedDomain "github.com/dutra/recodationsApp/internal/shared/domain"
	"github.com/dutra/recodationsApp/pkg/snowflake"
)

// InboxService manages transactional reception and idempotent processing of events.
type InboxService struct {
	repo repository.InboxRepository
}

func NewInboxService(repo repository.InboxRepository) *InboxService {
	return &InboxService{repo: repo}
}

// SaveEvent idempotently records an incoming event into the inbox table.
func (s *InboxService) SaveEvent(ctx context.Context, eventID string, eventType string, payload any) error {
	existing, err := s.repo.FindByEventID(ctx, eventID)
	if err == nil && existing != nil {
		return sharedDomain.ErrIdempotencyConflict
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed serializing inbox payload: %w", err)
	}

	now := time.Now().UTC()
	msg := &domain.InboxMessage{
		ID:        snowflake.Generate(),
		EventID:   eventID,
		EventType: eventType,
		Payload:   string(payloadBytes),
		Status:    domain.InboxStatusPending,
		Attempts:  0,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(ctx, msg)
}

// ProcessPendingEvents polls pending inbox messages and passes them to the registered handler.
func (s *InboxService) ProcessPendingEvents(ctx context.Context, limit int, handler func(ctx context.Context, event domain.InboxMessage) error) (int, error) {
	messages, err := s.repo.FindPending(ctx, limit)
	if err != nil {
		return 0, err
	}

	processedCount := 0
	for _, msg := range messages {
		slog.Info("Processing inbox message", "eventId", msg.EventID, "type", msg.EventType)

		hErr := handler(ctx, msg)
		if hErr != nil {
			slog.Error("Failed processing inbox event", "eventId", msg.EventID, "error", hErr)
			_ = msg.MarkFailed(hErr.Error())
		} else {
			_ = msg.MarkProcessed()
			processedCount++
		}

		if uErr := s.repo.Update(ctx, &msg); uErr != nil {
			slog.Error("Failed updating inbox message status", "eventId", msg.EventID, "error", uErr)
		}
	}

	return processedCount, nil
}
