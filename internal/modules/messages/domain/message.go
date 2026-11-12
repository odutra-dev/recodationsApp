package domain

import (
	"time"
)

// Message represents the core domain entity with state transitions and optimistic locking.
type Message struct {
	ID        int64         `json:"id" gorm:"primaryKey"`
	Name      string        `json:"name" gorm:"size:255;not null"`
	Content   string        `json:"content" gorm:"type:text;not null"`
	Image     *string       `json:"image" gorm:"type:text"`
	Status    MessageStatus `json:"status" gorm:"size:50;not null;default:'DRAFT'"`
	Version   int64         `json:"version" gorm:"default:1;not null"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// TransitionTo validates and performs the state transition following the Transition Pattern.
func (m *Message) TransitionTo(next MessageStatus) error {
	if err := MessageStateMachine.Validate(m.Status, next); err != nil {
		return err
	}
	m.Status = next
	return nil
}
