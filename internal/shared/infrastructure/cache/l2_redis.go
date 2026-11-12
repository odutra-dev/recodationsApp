package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements a distributed L2 cache using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis L2 cache instance.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string, dest any) error {
	if r.client == nil {
		return ErrCacheMiss
	}

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal(val, dest)
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if r.client == nil {
		return nil
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return nil
	}
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Clear(ctx context.Context) error {
	if r.client == nil {
		return nil
	}
	return r.client.FlushDB(ctx).Err()
}
