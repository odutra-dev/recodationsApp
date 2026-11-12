package domain

import (
	"time"

	"github.com/dutra/recodationsApp/internal/shared/domain"
)

type InboxStatus string

const (
	InboxStatusPending   InboxStatus = "PENDING"
	InboxStatusProcessed InboxStatus = "PROCESSED"
	InboxStatusFailed    InboxStatus = "FAILED"
)

// InboxStateMachine controls valid transitions of inbox events.
var InboxStateMachine = domain.NewStateMachine[InboxStatus]().
	Allow(InboxStatusPending, InboxStatusProcessed, InboxStatusFailed).
	Allow(InboxStatusFailed, InboxStatusPending) // Allows retrying failed events

// InboxMessage represents an idempotently stored incoming event / message.
type InboxMessage struct {
	ID           int64       `json:"id" gorm:"primaryKey"`
	EventID      string      `json:"eventId" gorm:"size:255;uniqueIndex;not null"`
	EventType    string      `json:"eventType" gorm:"size:100;not null"`
	Payload      string      `json:"payload" gorm:"type:text;not null"`
	Status       InboxStatus `json:"status" gorm:"size:50;not null;default:'PENDING'"`
	ErrorMessage *string     `json:"errorMessage" gorm:"type:text"`
	Attempts     int         `json:"attempts" gorm:"default:0;not null"`
	ProcessedAt  *time.Time  `json:"processedAt"`
	Version      int64       `json:"version" gorm:"default:1;not null"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

// MarkProcessed transitions the status to PROCESSED.
func (im *InboxMessage) MarkProcessed() error {
	if err := InboxStateMachine.Validate(im.Status, InboxStatusProcessed); err != nil {
		return err
	}
	now := time.Now().UTC()
	im.Status = InboxStatusProcessed
	im.ProcessedAt = &now
	im.ErrorMessage = nil
	return nil
}

// MarkFailed transitions the status to FAILED and increments attempts.
func (im *InboxMessage) MarkFailed(errMsg string) error {
	if err := InboxStateMachine.Validate(im.Status, InboxStatusFailed); err != nil {
		return err
	}
	im.Status = InboxStatusFailed
	im.ErrorMessage = &errMsg
	im.Attempts++
	return nil
}
