package cache

import (
	"context"
	"log/slog"
	"time"
)

// TieredCache combines an L1 in-memory cache and an L2 distributed (Redis) cache.
type TieredCache struct {
	l1 Cache
	l2 Cache
}

// NewTieredCache constructs a two-level cache.
func NewTieredCache(l1 Cache, l2 Cache) *TieredCache {
	return &TieredCache{
		l1: l1,
		l2: l2,
	}
}

// Get checks L1 first. If missed, checks L2 and backfills L1 upon hit.
func (t *TieredCache) Get(ctx context.Context, key string, dest any) error {
	// 1. Check L1 Cache
	if t.l1 != nil {
		err := t.l1.Get(ctx, key, dest)
		if err == nil {
			return nil
		}
	}

	// 2. Check L2 Cache
	if t.l2 != nil {
		err := t.l2.Get(ctx, key, dest)
		if err == nil {
			// Backfill L1 with default short TTL (e.g. 5 minutes)
			if t.l1 != nil {
				_ = t.l1.Set(ctx, key, dest, 5*time.Minute)
			}
			return nil
		}
		if err != ErrCacheMiss {
			slog.Warn("L2 cache read error (continuing as miss)", "error", err, "key", key)
		}
	}

	return ErrCacheMiss
}

// Set writes to both L1 and L2 caches.
func (t *TieredCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	var firstErr error

	if t.l1 != nil {
		if err := t.l1.Set(ctx, key, value, ttl); err != nil {
			firstErr = err
		}
	}

	if t.l2 != nil {
		if err := t.l2.Set(ctx, key, value, ttl); err != nil {
			slog.Warn("L2 cache write error", "error", err, "key", key)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

// Delete removes the key from both cache layers.
func (t *TieredCache) Delete(ctx context.Context, key string) error {
	if t.l1 != nil {
		_ = t.l1.Delete(ctx, key)
	}
	if t.l2 != nil {
		_ = t.l2.Delete(ctx, key)
	}
	return nil
}

// Clear flushes both cache layers.
func (t *TieredCache) Clear(ctx context.Context) error {
	if t.l1 != nil {
		_ = t.l1.Clear(ctx)
	}
	if t.l2 != nil {
		_ = t.l2.Clear(ctx)
	}
	return nil
}
