package media

import (
	"context"
	"errors"
	"io"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/pkg/storage"
)

// Objects is the MinIO adapter for primary-image bytes.
type Objects struct {
	store storage.Store
}

// NewObjects binds primary-image put/get/delete to an object store.
func NewObjects(store storage.Store) *Objects {
	return &Objects{store: store}
}

// Put stores or replaces an object.
func (o *Objects) Put(ctx context.Context, key, contentType string, body io.Reader, size int64) error {
	if o == nil || o.store == nil {
		return errors.New("primary images are not configured")
	}
	return o.store.Put(ctx, key, contentType, body, size)
}

// Get opens an object. Missing keys are not found.
func (o *Objects) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	if o == nil || o.store == nil {
		return nil, "", errors.New("primary images are not configured")
	}
	body, contentType, err := o.store.Get(ctx, key)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, "", apperrors.ErrNotFound
	}
	return body, contentType, err
}

// Delete removes an object. Missing keys succeed.
func (o *Objects) Delete(ctx context.Context, key string) error {
	if o == nil || o.store == nil {
		return errors.New("primary images are not configured")
	}
	return o.store.Delete(ctx, key)
}
