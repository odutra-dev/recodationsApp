package dto

import "time"

// CreateMessageInput represents input to create a message.
type CreateMessageInput struct {
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Image   *string `json:"image,omitempty"`
}

// UpdateMessageInput represents input to update message fields.
type UpdateMessageInput struct {
	Name    *string `json:"name,omitempty"`
	Content *string `json:"content,omitempty"`
	Image   *string `json:"image,omitempty"`
	Version int64   `json:"version"` // Required for Optimistic Locking
}

// TransitionMessageInput represents input to change the message state.
type TransitionMessageInput struct {
	Status  string `json:"status"`
	Version int64  `json:"version"` // Required for Optimistic Locking
}

// MessageResponse represents the public payload for a message.
type MessageResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Image     *string   `json:"image"`
	Status    string    `json:"status"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
