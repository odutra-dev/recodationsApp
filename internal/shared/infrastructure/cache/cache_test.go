package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sampleData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func TestMemoryCache(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache()

	key := "test:item:1"
	item := sampleData{ID: 100, Name: "Golang"}

	// Miss
	var dest sampleData
	err := c.Get(ctx, key, &dest)
	assert.ErrorIs(t, err, ErrCacheMiss)

	// Set & Hit
	err = c.Set(ctx, key, item, time.Minute)
	require.NoError(t, err)

	err = c.Get(ctx, key, &dest)
	require.NoError(t, err)
	assert.Equal(t, item, dest)

	// Delete
	err = c.Delete(ctx, key)
	require.NoError(t, err)

	err = c.Get(ctx, key, &dest)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestMemoryCacheExpiration(t *testing.T) {
	ctx := context.Background()
	c := NewMemoryCache()

	key := "test:item:expire"
	item := sampleData{ID: 200, Name: "Expire"}

	err := c.Set(ctx, key, item, 20*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(35 * time.Millisecond)

	var dest sampleData
	err = c.Get(ctx, key, &dest)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestTieredCacheFallback(t *testing.T) {
	ctx := context.Background()
	l1 := NewMemoryCache()
	// L2 is nil or mock
	tiered := NewTieredCache(l1, nil)

	key := "tiered:1"
	item := sampleData{ID: 300, Name: "Tiered"}

	err := tiered.Set(ctx, key, item, time.Minute)
	require.NoError(t, err)

	var dest sampleData
	err = tiered.Get(ctx, key, &dest)
	require.NoError(t, err)
	assert.Equal(t, item, dest)
}
