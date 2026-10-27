package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"inventory_dashboard/internal/models"
	"inventory_dashboard/pkg/logger"
	"inventory_dashboard/pkg/metrics"
)

const (
	keyMakes = "catalog:makes"
)

// Reader loads makes and models from PostgreSQL (primary or replica).
type Reader interface {
	ListMakes(ctx context.Context) ([]models.Make, error)
	ListModels(ctx context.Context, makeID *int32) ([]models.Model, error)
}

// Catalog implements the catalog service.
type Catalog struct {
	repo  Reader
	cache Cache
	ttl   time.Duration
}

// New binds catalog rules to a catalog reader.
// cache may be nil (always read Postgres). ttl defaults to 1h when zero.
func New(repo Reader, cache Cache, ttl time.Duration) *Catalog {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &Catalog{repo: repo, cache: cache, ttl: ttl}
}

// ListMakes returns the seeded makes.
func (c *Catalog) ListMakes(ctx context.Context) ([]models.Make, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("catalog is not configured")
	}
	if raw, ok := c.cacheGet(ctx, keyMakes, "makes"); ok {
		var makes []models.Make
		if err := json.Unmarshal(raw, &makes); err == nil {
			logger.FromContext(ctx).Debug("listed makes",
				zap.String("layer", "logic"),
				zap.String("op", "ListMakes"),
				zap.String("cache", "hit"),
				zap.Int("count", len(makes)),
			)
			return makes, nil
		}
		c.cacheError("makes", "unmarshal")
	}

	makes, err := c.repo.ListMakes(ctx)
	if err != nil {
		return nil, err
	}
	c.cacheSet(ctx, keyMakes, "makes", makes)
	logger.FromContext(ctx).Debug("listed makes",
		zap.String("layer", "logic"),
		zap.String("op", "ListMakes"),
		zap.String("cache", "miss"),
		zap.Int("count", len(makes)),
	)
	return makes, nil
}

// ListModels returns seeded models. An empty makeID returns every model.
func (c *Catalog) ListModels(ctx context.Context, makeID string) ([]models.Model, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("catalog is not configured")
	}
	parsed, err := ParseMakeID(makeID)
	if err != nil {
		return nil, err
	}
	key := modelsKey(parsed)
	if raw, ok := c.cacheGet(ctx, key, "models"); ok {
		var list []models.Model
		if err := json.Unmarshal(raw, &list); err == nil {
			logger.FromContext(ctx).Debug("listed models",
				zap.String("layer", "logic"),
				zap.String("op", "ListModels"),
				zap.String("cache", "hit"),
				zap.Int("count", len(list)),
			)
			return list, nil
		}
		c.cacheError("models", "unmarshal")
	}

	list, err := c.repo.ListModels(ctx, parsed)
	if err != nil {
		return nil, err
	}
	c.cacheSet(ctx, key, "models", list)
	logger.FromContext(ctx).Debug("listed models",
		zap.String("layer", "logic"),
		zap.String("op", "ListModels"),
		zap.String("cache", "miss"),
		zap.Int("count", len(list)),
	)
	return list, nil
}

func modelsKey(makeID *int32) string {
	if makeID == nil {
		return "catalog:models"
	}
	return "catalog:models:make:" + strconv.FormatInt(int64(*makeID), 10)
}

func (c *Catalog) cacheGet(ctx context.Context, key, resource string) ([]byte, bool) {
	if c.cache == nil {
		metrics.CatalogCacheMisses.WithLabelValues(resource).Inc()
		return nil, false
	}
	raw, hit, err := c.cache.Get(ctx, key)
	if err != nil {
		c.cacheError(resource, "get")
		logger.FromContext(ctx).Warn("catalog cache get failed",
			zap.String("resource", resource),
			zap.Error(err),
		)
		metrics.CatalogCacheMisses.WithLabelValues(resource).Inc()
		return nil, false
	}
	if !hit {
		metrics.CatalogCacheMisses.WithLabelValues(resource).Inc()
		return nil, false
	}
	metrics.CatalogCacheHits.WithLabelValues(resource).Inc()
	return raw, true
}

func (c *Catalog) cacheSet(ctx context.Context, key, resource string, value any) {
	if c.cache == nil {
		return
	}
	raw, err := json.Marshal(value)
	if err != nil {
		c.cacheError(resource, "marshal")
		return
	}
	if err := c.cache.Set(ctx, key, raw, c.ttl); err != nil {
		c.cacheError(resource, "set")
		logger.FromContext(ctx).Warn("catalog cache set failed",
			zap.String("resource", resource),
			zap.Error(err),
		)
	}
}

func (c *Catalog) cacheError(resource, op string) {
	metrics.CatalogCacheErrors.WithLabelValues(resource, op).Inc()
}
