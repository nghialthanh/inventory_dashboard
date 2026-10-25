package services

import (
	"context"
	"fmt"

	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/models"
	"inventory_dashboard/internal/repositories"
)

// Services is the application registry. The process constructs one in Run.
type Services struct {
	catalog *repositories.Catalog
}

// New returns the application services. A nil catalog leaves catalog reads unavailable.
func New(catalog *repositories.Catalog) *Services {
	return &Services{catalog: catalog}
}

// Health reports that the process is up. It does not check PostgreSQL or MinIO.
func (s *Services) Health() error {
	return nil
}

// ListMakes returns the seeded makes.
func (s *Services) ListMakes(ctx context.Context) ([]models.Make, error) {
	if s.catalog == nil {
		return nil, fmt.Errorf("catalog is not configured")
	}
	return s.catalog.ListMakes(ctx)
}

// ListModels returns seeded models. An empty makeID returns every model.
func (s *Services) ListModels(ctx context.Context, makeID string) ([]models.Model, error) {
	if s.catalog == nil {
		return nil, fmt.Errorf("catalog is not configured")
	}
	parsed, err := catalog.ParseMakeID(makeID)
	if err != nil {
		return nil, err
	}
	return s.catalog.ListModels(ctx, parsed)
}
