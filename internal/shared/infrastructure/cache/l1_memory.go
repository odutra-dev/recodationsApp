package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type memoryItem struct {
	data      []byte
	expiresAt time.Time
}

// MemoryCache implements an in-memory L1 cache with TTL support and thread safety.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

// NewMemoryCache creates a new in-memory cache instance.
func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		items: make(map[string]memoryItem),
	}
	// Background cleanup of expired keys every 60 seconds
	go mc.startJanitor(60 * time.Second)
	return mc
}

func (m *MemoryCache) Get(ctx context.Context, key string, dest any) error {
	m.mu.RLock()
	item, found := m.items[key]
	m.mu.RUnlock()

	if !found {
		return ErrCacheMiss
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		m.Delete(ctx, key)
		return ErrCacheMiss
	}

	return json.Unmarshal(item.data, dest)
}

func (m *MemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	m.mu.Lock()
	m.items[key] = memoryItem{
		data:      bytes,
		expiresAt: expiresAt,
	}
	m.mu.Unlock()

	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) Clear(ctx context.Context) error {
	m.mu.Lock()
	m.items = make(map[string]memoryItem)
	m.mu.Unlock()
	return nil
}

func (m *MemoryCache) startJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		now := time.Now()
		m.mu.Lock()
		for k, v := range m.items {
			if !v.expiresAt.IsZero() && now.After(v.expiresAt) {
				delete(m.items, k)
			}
		}
		m.mu.Unlock()
	}
}
