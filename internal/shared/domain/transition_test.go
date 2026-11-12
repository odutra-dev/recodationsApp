package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStatus string

const (
	StatusDraft     TestStatus = "DRAFT"
	StatusPublished TestStatus = "PUBLISHED"
	StatusArchived  TestStatus = "ARCHIVED"
)

func TestStateMachine(t *testing.T) {
	sm := NewStateMachine[TestStatus]().
		Allow(StatusDraft, StatusPublished, StatusArchived).
		Allow(StatusPublished, StatusArchived, StatusDraft).
		Allow(StatusArchived) // Terminal state: cannot transition anywhere

	// Valid transitions
	assert.True(t, sm.CanTransition(StatusDraft, StatusPublished))
	assert.True(t, sm.CanTransition(StatusDraft, StatusArchived))
	assert.True(t, sm.CanTransition(StatusPublished, StatusArchived))
	assert.True(t, sm.CanTransition(StatusPublished, StatusDraft))
	require.NoError(t, sm.Validate(StatusDraft, StatusPublished))

	// Invalid transitions
	assert.False(t, sm.CanTransition(StatusArchived, StatusDraft))
	assert.False(t, sm.CanTransition(StatusArchived, StatusPublished))
	err := sm.Validate(StatusArchived, StatusPublished)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidStateTransition)
}
