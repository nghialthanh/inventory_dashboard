package catalog

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"inventory_dashboard/internal/models"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/logger"
)

// Catalog implements the catalog service.
type Catalog struct {
	repo *repositories.Catalog
}

// New binds catalog rules to the catalog repository.
func New(repo *repositories.Catalog) *Catalog {
	return &Catalog{repo: repo}
}

// ListMakes returns the seeded makes.
func (c *Catalog) ListMakes(ctx context.Context) ([]models.Make, error) {
	if c.repo == nil {
		return nil, fmt.Errorf("catalog is not configured")
	}
	makes, err := c.repo.ListMakes(ctx)
	if err != nil {
		return nil, err
	}
	logger.FromContext(ctx).Debug("listed makes",
		zap.String("layer", "logic"),
		zap.String("op", "ListMakes"),
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
	list, err := c.repo.ListModels(ctx, parsed)
	if err != nil {
		return nil, err
	}
	logger.FromContext(ctx).Debug("listed models",
		zap.String("layer", "logic"),
		zap.String("op", "ListModels"),
		zap.Int("count", len(list)),
	)
	return list, nil
}
