package catalog

import (
	"context"
	"time"
)

// Cache stores serialized catalog list payloads. Implementations may fail;
// callers fail open to PostgreSQL.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}
