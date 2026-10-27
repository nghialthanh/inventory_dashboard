package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis is a thin byte cache over Redis. Callers treat errors as miss/fail-open.
type Redis struct {
	client *redis.Client
}

// Open connects to Redis. redisURL is a redis:// URL.
func Open(redisURL string) (*Redis, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis url: %w", err)
	}
	client := redis.NewClient(opts)
	return &Redis{client: client}, nil
}

// Ping checks connectivity.
func (r *Redis) Ping(ctx context.Context) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("redis is not configured")
	}
	return r.client.Ping(ctx).Err()
}

// Get returns the value and whether it was a hit.
func (r *Redis) Get(ctx context.Context, key string) ([]byte, bool, error) {
	if r == nil || r.client == nil {
		return nil, false, nil
	}
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return val, true, nil
}

// Set stores value with TTL.
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Close closes the client.
func (r *Redis) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}
