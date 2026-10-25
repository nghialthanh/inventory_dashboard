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
func (c *Catalog) ModelMakeID(ctx context.Context, modelID string) (string, bool, error) {
	if c == nil || c.pool == nil {
		return "", false, errors.New("catalog is not configured")
	}
	var makeID string
	err := c.pool.QueryRow(ctx, `SELECT make_id::text FROM models WHERE id = $1`, modelID).Scan(&makeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return makeID, true, nil
}

// ListMakes returns every make, ordered by name.
func (c *Catalog) ListMakes(ctx context.Context) ([]models.Make, error) {
	rows, err := c.pool.Query(ctx, `SELECT id::text, name FROM makes ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	makes := []models.Make{}
	for rows.Next() {
		var make models.Make
		if err := rows.Scan(&make.ID, &make.Name); err != nil {
			return nil, err
		}
		makes = append(makes, make)
	}
	return makes, rows.Err()
}

// ListModels returns models, optionally limited to one make, ordered by name.
func (c *Catalog) ListModels(ctx context.Context, makeID *string) ([]models.Model, error) {
	query := `SELECT id::text, make_id::text, name FROM models`
	args := []any{}
	if makeID != nil {
		query += ` WHERE make_id = $1`
		args = append(args, *makeID)
	}
	query += ` ORDER BY name`

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []models.Model{}
	for rows.Next() {
		var model models.Model
		if err := rows.Scan(&model.ID, &model.MakeID, &model.Name); err != nil {
			return nil, err
		}
		list = append(list, model)
	}
	return list, rows.Err()
}
