package storage

import (
	"bytes"
	"context"
	"io"
	"sync"
)

type memoryObject struct {
	contentType string
	data        []byte
}

// Memory is an in-process object store for HTTP tests.
type Memory struct {
	mu      sync.Mutex
	objects map[string]memoryObject
}

// NewMemory builds an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{objects: map[string]memoryObject{}}
}

// Put stores or replaces an object.
func (m *Memory) Put(_ context.Context, key, contentType string, body io.Reader, _ int64) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.objects[key] = memoryObject{contentType: contentType, data: data}
	return nil
}

// Get opens an object.
func (m *Memory) Get(_ context.Context, key string) (io.ReadCloser, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	obj, ok := m.objects[key]
	if !ok {
		return nil, "", ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(obj.data)), obj.contentType, nil
}

// Delete removes an object. Missing keys succeed.
func (m *Memory) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.objects, key)
	return nil
}
