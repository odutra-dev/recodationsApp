package domain

import "github.com/dutra/recodationsApp/internal/shared/domain"

type MessageStatus string

const (
	MessageStatusDraft     MessageStatus = "DRAFT"
	MessageStatusPublished MessageStatus = "PUBLISHED"
	MessageStatusArchived  MessageStatus = "ARCHIVED"
)

// MessageStateMachine defines the allowed state transitions for messages.
var MessageStateMachine = domain.NewStateMachine[MessageStatus]().
	Allow(MessageStatusDraft, MessageStatusPublished, MessageStatusArchived).
	Allow(MessageStatusPublished, MessageStatusArchived, MessageStatusDraft).
	Allow(MessageStatusArchived) // Terminal state
