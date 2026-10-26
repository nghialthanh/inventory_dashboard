package storage

import (
	"context"
	"errors"
	"io"
)

// ErrNotFound means the object key is missing from the bucket.
var ErrNotFound = errors.New("object not found")

// Store holds primary-image object bytes. Production uses MinIO; tests use memory.
type Store interface {
	Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error
	Get(ctx context.Context, key string) (body io.ReadCloser, contentType string, err error)
	Delete(ctx context.Context, key string) error
}
