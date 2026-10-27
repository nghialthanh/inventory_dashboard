package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/common/snowflake"
	"inventory_dashboard/internal/models"
)

// Vehicles reads and records vehicles.
type Vehicles struct {
	pool *pgxpool.Pool
}

// NewVehicles binds vehicle commands to PostgreSQL.
func NewVehicles(pool *pgxpool.Pool) *Vehicles {
	return &Vehicles{pool: pool}
}

// Insert records an in-stock vehicle with a pre-minted address. A duplicate VIN is a conflict.
func (v *Vehicles) Insert(ctx context.Context, id int64, createdAt time.Time, vin string, modelID int32, enteredOn time.Time) error {
	tx, err := v.pool.Begin(ctx)
	if err != nil {
		return logRepoErr(ctx, "vehicles.Insert.begin", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO vehicle_vins (vin, vehicle_id)
		VALUES ($1, $2)
	`, vin, id)
	if err != nil {
		return logRepoErr(ctx, "vehicles.Insert.vin", vinConflict(err))
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO vehicles (id, vin, model_id, entered_on, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
	`, id, vin, modelID, enteredOn, createdAt.UTC())
	if err != nil {
		return logRepoErr(ctx, "vehicles.Insert.row", vinConflict(err))
	}
	if err := tx.Commit(ctx); err != nil {
		return logRepoErr(ctx, "vehicles.Insert.commit", err)
	}
	return nil
}

// Correct changes the model or entered on of one vehicle.
func (v *Vehicles) Correct(ctx context.Context, id int64, modelID *int32, enteredOn *time.Time) error {
	start, end := snowflake.MonthBounds(id)
	tag, err := v.pool.Exec(ctx, `
		UPDATE vehicles
		SET model_id = COALESCE($2, model_id),
			entered_on = COALESCE($3, entered_on),
			updated_at = now()
		WHERE id = $1
			AND created_at >= $4 AND created_at < $5
			AND stocked_out_on IS NULL
	`, id, modelID, enteredOn, start, end)
	if err != nil {
		return logRepoErr(ctx, "vehicles.Correct", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// StockOut records or replaces stocked-out on for one vehicle.
func (v *Vehicles) StockOut(ctx context.Context, id int64, stockedOutOn time.Time) error {
	start, end := snowflake.MonthBounds(id)
	tag, err := v.pool.Exec(ctx, `
		UPDATE vehicles
		SET stocked_out_on = $2,
			updated_at = now()
		WHERE id = $1
			AND created_at >= $3 AND created_at < $4
	`, id, stockedOutOn, start, end)
	if err != nil {
		return logRepoErr(ctx, "vehicles.StockOut", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

// SetProposedAction stores or clears the proposed action on one vehicle.
func (v *Vehicles) SetProposedAction(ctx context.Context, id int64, proposedAction *string) error {
	start, end := snowflake.MonthBounds(id)
	tag, err := v.pool.Exec(ctx, `
		UPDATE vehicles
		SET proposed_action = $2,
			updated_at = now()
		WHERE id = $1
			AND created_at >= $3 AND created_at < $4
	`, id, proposedAction, start, end)
	if err != nil {
		return logRepoErr(ctx, "vehicles.SetProposedAction", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func vinConflict(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return &apperrors.WithMessage{Err: apperrors.ErrConflict, Message: "VIN is already used"}
	}
	return err
}

const vehicleSelect = `
	SELECT v.id, v.vin,
		mk.id, mk.name,
		md.id, md.name,
		v.entered_on, v.stocked_out_on, v.proposed_action
	FROM vehicles v
	JOIN models md ON md.id = v.model_id
	JOIN makes mk ON mk.id = md.make_id
`

// Get loads one vehicle by its address, pruning to the snowflake month.
func (v *Vehicles) Get(ctx context.Context, id int64) (models.StoredVehicle, error) {
	start, end := snowflake.MonthBounds(id)
	row := v.pool.QueryRow(ctx, vehicleSelect+`
		WHERE v.id = $1
			AND v.created_at >= $2 AND v.created_at < $3
	`, id, start, end)
	vehicle, err := scanVehicle(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.StoredVehicle{}, apperrors.ErrNotFound
	}
	if err != nil {
		return models.StoredVehicle{}, logRepoErr(ctx, "vehicles.Get", err)
	}
	return vehicle, nil
}

// ListFilter narrows the inventory list.
type ListFilter struct {
	MakeID       int32
	ModelID      int32
	StockState   string
	StockedOutOn string
	// Aging is nil when omitted, true for aging stock only, false for not aging stock.
	Aging      *bool
	AgingBound string // calendar day: in stock entered before this day is aging stock
	Limit      int
	Offset     int
}

// List returns a page of vehicles ordered by entered on, then id, and the total.
func (v *Vehicles) List(ctx context.Context, filter ListFilter) ([]models.StoredVehicle, int, error) {
	where, args := filterSQL(filter)
	var total int
	if err := v.pool.QueryRow(ctx, `
		SELECT count(*)
		FROM vehicles v
		JOIN models md ON md.id = v.model_id
		JOIN makes mk ON mk.id = md.make_id
	`+where, args...).Scan(&total); err != nil {
		return nil, 0, logRepoErr(ctx, "vehicles.List.count", err)
	}

	listArgs := append(append([]any{}, args...), filter.Limit, filter.Offset)
	rows, err := v.pool.Query(ctx, vehicleSelect+where+fmt.Sprintf(`
		ORDER BY v.entered_on ASC, v.id ASC
		LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2), listArgs...)
	if err != nil {
		return nil, 0, logRepoErr(ctx, "vehicles.List.query", err)
	}
	defer rows.Close()

	vehicles := []models.StoredVehicle{}
	for rows.Next() {
		vehicle, err := scanVehicle(rows)
		if err != nil {
			return nil, 0, logRepoErr(ctx, "vehicles.List.scan", err)
		}
		vehicles = append(vehicles, vehicle)
	}
	return vehicles, total, logRepoErr(ctx, "vehicles.List.rows", rows.Err())
}

func filterSQL(filter ListFilter) (string, []any) {
	where := " WHERE 1=1"
	args := []any{}
	if filter.MakeID != 0 {
		args = append(args, filter.MakeID)
		where += fmt.Sprintf(" AND mk.id = $%d", len(args))
	}
	if filter.ModelID != 0 {
		args = append(args, filter.ModelID)
		where += fmt.Sprintf(" AND md.id = $%d", len(args))
	}
	switch filter.StockState {
	case "in_stock":
		where += " AND v.stocked_out_on IS NULL"
	case "stocked_out":
		where += " AND v.stocked_out_on IS NOT NULL"
	}
	if filter.StockedOutOn != "" {
		args = append(args, filter.StockedOutOn)
		where += fmt.Sprintf(" AND v.stocked_out_on = $%d::date", len(args))
	}
	if filter.Aging != nil {
		args = append(args, filter.AgingBound)
		if *filter.Aging {
			where += fmt.Sprintf(
				" AND v.stocked_out_on IS NULL AND v.entered_on < $%d::date",
				len(args),
			)
		} else {
			where += fmt.Sprintf(
				" AND (v.stocked_out_on IS NOT NULL OR v.entered_on >= $%d::date)",
				len(args),
			)
		}
	}
	return where, args
}

type scanner interface {
	Scan(dest ...any) error
}

func scanVehicle(row scanner) (models.StoredVehicle, error) {
	var vehicle models.StoredVehicle
	err := row.Scan(
		&vehicle.ID,
		&vehicle.VIN,
		&vehicle.Make.ID,
		&vehicle.Make.Name,
		&vehicle.Model.ID,
		&vehicle.Model.Name,
		&vehicle.EnteredOn,
		&vehicle.StockedOutOn,
		&vehicle.ProposedAction,
	)
	if err != nil {
		return vehicle, err
	}
	return vehicle, nil
}
