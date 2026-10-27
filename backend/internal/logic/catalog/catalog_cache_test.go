package catalog_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/models"
)

type memCache struct {
	data   map[string][]byte
	getErr error
	setErr error
	gets   int
	sets   int
}

func (m *memCache) Get(_ context.Context, key string) ([]byte, bool, error) {
	m.gets++
	if m.getErr != nil {
		return nil, false, m.getErr
	}
	v, ok := m.data[key]
	return v, ok, nil
}

func (m *memCache) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	m.sets++
	if m.setErr != nil {
		return m.setErr
	}
	if m.data == nil {
		m.data = map[string][]byte{}
	}
	m.data[key] = append([]byte(nil), value...)
	return nil
}

type stubReader struct {
	makes     []models.Make
	models    []models.Model
	listCalls int
}

func (s *stubReader) ListMakes(context.Context) ([]models.Make, error) {
	s.listCalls++
	return s.makes, nil
}

func (s *stubReader) ListModels(context.Context, *int32) ([]models.Model, error) {
	s.listCalls++
	return s.models, nil
}

func TestListMakesCacheHitSkipsReader(t *testing.T) {
	payload, err := json.Marshal([]models.Make{{ID: 1, Name: "Ford"}})
	if err != nil {
		t.Fatal(err)
	}
	c := &memCache{data: map[string][]byte{"catalog:makes": payload}}
	r := &stubReader{makes: []models.Make{{ID: 9, Name: "Other"}}}
	svc := catalog.New(r, c, time.Minute)

	got, err := svc.ListMakes(context.Background())
	if err != nil {
		t.Fatalf("ListMakes: %v", err)
	}
	if r.listCalls != 0 {
		t.Fatalf("reader called %d times, want 0 on cache hit", r.listCalls)
	}
	if len(got) != 1 || got[0].Name != "Ford" {
		t.Fatalf("got %#v, want Ford from cache", got)
	}
}

func TestListMakesCacheMissFillsCache(t *testing.T) {
	c := &memCache{}
	r := &stubReader{makes: []models.Make{{ID: 1, Name: "Honda"}}}
	svc := catalog.New(r, c, time.Minute)

	got, err := svc.ListMakes(context.Background())
	if err != nil {
		t.Fatalf("ListMakes: %v", err)
	}
	if r.listCalls != 1 {
		t.Fatalf("reader calls = %d, want 1", r.listCalls)
	}
	if len(got) != 1 || got[0].Name != "Honda" {
		t.Fatalf("got %#v", got)
	}
	if c.sets != 1 || c.data["catalog:makes"] == nil {
		t.Fatalf("expected cache set for catalog:makes")
	}
}

func TestListMakesFailOpenOnCacheGetError(t *testing.T) {
	c := &memCache{getErr: errors.New("redis down")}
	r := &stubReader{makes: []models.Make{{ID: 2, Name: "Toyota"}}}
	svc := catalog.New(r, c, time.Minute)

	got, err := svc.ListMakes(context.Background())
	if err != nil {
		t.Fatalf("ListMakes should fail open: %v", err)
	}
	if r.listCalls != 1 {
		t.Fatalf("reader calls = %d, want 1 after cache error", r.listCalls)
	}
	if len(got) != 1 || got[0].Name != "Toyota" {
		t.Fatalf("got %#v", got)
	}
}

func TestListModelsKeysByMake(t *testing.T) {
	c := &memCache{}
	r := &stubReader{models: []models.Model{{ID: 3, MakeID: 2, Name: "Accord"}}}
	svc := catalog.New(r, c, time.Minute)

	if _, err := svc.ListModels(context.Background(), "2"); err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	if c.data["catalog:models:make:2"] == nil {
		t.Fatalf("missing keyed cache entry, keys=%v", keysOf(c.data))
	}

	c2 := &memCache{}
	r2 := &stubReader{models: []models.Model{{ID: 1, MakeID: 1, Name: "F-150"}}}
	svc2 := catalog.New(r2, c2, time.Minute)
	if _, err := svc2.ListModels(context.Background(), ""); err != nil {
		t.Fatalf("ListModels all: %v", err)
	}
	if c2.data["catalog:models"] == nil {
		t.Fatalf("missing all-models key, keys=%v", keysOf(c2.data))
	}
}

func keysOf(m map[string][]byte) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
