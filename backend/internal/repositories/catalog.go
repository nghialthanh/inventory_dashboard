package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"inventory_dashboard/internal/models"
)

// Catalog reads the seeded make and model catalog.
type Catalog struct {
	pool *pgxpool.Pool
}

// NewCatalog binds the catalog to PostgreSQL.
func NewCatalog(pool *pgxpool.Pool) *Catalog {
	return &Catalog{pool: pool}
}

// ModelMakeID returns the make that owns the model. ok is false when the model is unknown.
func (c *Catalog) ModelMakeID(ctx context.Context, modelID int32) (int32, bool, error) {
	if c == nil || c.pool == nil {
		return 0, false, errors.New("catalog is not configured")
	}
	var makeID int32
	err := c.pool.QueryRow(ctx, `SELECT make_id FROM models WHERE id = $1`, modelID).Scan(&makeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, logRepoErr(ctx, "catalog.ModelMakeID", err)
	}
	return makeID, true, nil
}

// ListMakes returns every make, ordered by name.
func (c *Catalog) ListMakes(ctx context.Context) ([]models.Make, error) {
	rows, err := c.pool.Query(ctx, `SELECT id, name FROM makes ORDER BY name`)
	if err != nil {
		return nil, logRepoErr(ctx, "catalog.ListMakes", err)
	}
	defer rows.Close()

	makes := []models.Make{}
	for rows.Next() {
		var make models.Make
		if err := rows.Scan(&make.ID, &make.Name); err != nil {
			return nil, logRepoErr(ctx, "catalog.ListMakes.scan", err)
		}
		makes = append(makes, make)
	}
	return makes, logRepoErr(ctx, "catalog.ListMakes.rows", rows.Err())
}

// ListModels returns models, optionally limited to one make, ordered by name.
func (c *Catalog) ListModels(ctx context.Context, makeID *int32) ([]models.Model, error) {
	query := `SELECT id, make_id, name FROM models`
	args := []any{}
	if makeID != nil {
		query += ` WHERE make_id = $1`
		args = append(args, *makeID)
	}
	query += ` ORDER BY name`

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, logRepoErr(ctx, "catalog.ListModels", err)
	}
	defer rows.Close()

	list := []models.Model{}
	for rows.Next() {
		var model models.Model
		if err := rows.Scan(&model.ID, &model.MakeID, &model.Name); err != nil {
			return nil, logRepoErr(ctx, "catalog.ListModels.scan", err)
		}
		list = append(list, model)
	}
	return list, logRepoErr(ctx, "catalog.ListModels.rows", rows.Err())
}
