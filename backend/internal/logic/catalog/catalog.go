package catalog

import (
	"context"
	"fmt"

	"inventory_dashboard/internal/models"
	"inventory_dashboard/internal/repositories"
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
	return c.repo.ListMakes(ctx)
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
	return c.repo.ListModels(ctx, parsed)
}
