package services

import (
	"context"

	"inventory_dashboard/internal/models"
)

// Catalog is the read-only make and model catalog.
type Catalog interface {
	ListMakes(ctx context.Context) ([]models.Make, error)
	ListModels(ctx context.Context, makeID string) ([]models.Model, error)
}
