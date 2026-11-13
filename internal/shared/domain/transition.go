package domain

import (
	"fmt"
	"time"
)

// StateTransition represents a historical log of a state change.
type StateTransition[T comparable] struct {
	From      T         `json:"from"`
	To        T         `json:"to"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
}

// StateMachine manages allowed state transitions for a given domain entity state.
type StateMachine[T comparable] struct {
	allowed map[T][]T
}

// NewStateMachine creates a new state machine configuration.
func NewStateMachine[T comparable]() *StateMachine[T] {
	return &StateMachine[T]{
		allowed: make(map[T][]T),
	}
}

// Allow registers permitted destination states from a specific source state.
func (sm *StateMachine[T]) Allow(from T, to ...T) *StateMachine[T] {
	sm.allowed[from] = append(sm.allowed[from], to...)
	return sm
}

// CanTransition returns true if the transition from current to next is permitted.
func (sm *StateMachine[T]) CanTransition(from, to T) bool {
	destinations, exists := sm.allowed[from]
	if !exists {
		return false
	}
	for _, d := range destinations {
		if d == to {
			return true
		}
	}
	return false
}

// Validate verifies if transitioning from 'from' to 'to' is allowed; returns ErrInvalidStateTransition if not.
func (sm *StateMachine[T]) Validate(from, to T) error {
	if !sm.CanTransition(from, to) {
		return fmt.Errorf("%w: cannot transition from %v to %v", ErrInvalidStateTransition, from, to)
	}
	return nil
}
