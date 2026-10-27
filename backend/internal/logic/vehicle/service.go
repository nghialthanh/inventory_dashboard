package vehicle

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"inventory_dashboard/internal/common/constant"
	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/common/snowflake"
	"inventory_dashboard/internal/logic/aging"
	"inventory_dashboard/internal/models"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/logger"
)

// Vehicles implements the vehicle service.
type Vehicles struct {
	repo    *repositories.Vehicles
	catalog *repositories.Catalog
	ids     *snowflake.Generator
	today   func() time.Time
}

// New binds vehicle rules to the vehicle and catalog repositories.
// A nil today function uses the process clock.
func New(repo *repositories.Vehicles, catalog *repositories.Catalog, ids *snowflake.Generator, today func() time.Time) *Vehicles {
	if today == nil {
		today = time.Now
	}
	return &Vehicles{repo: repo, catalog: catalog, ids: ids, today: today}
}

// RecordVehicle stores an in-stock vehicle and returns it with derived fields.
func (v *Vehicles) RecordVehicle(ctx context.Context, vin string, modelID int32, enteredOn string) (models.Vehicle, error) {
	if v.repo == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	normalized, err := NormalizeVIN(vin)
	if err != nil {
		return models.Vehicle{}, err
	}
	if err := AcceptModelID(modelID); err != nil {
		return models.Vehicle{}, err
	}
	entered, err := ParseEnteredOn(enteredOn, v.today())
	if err != nil {
		return models.Vehicle{}, err
	}
	if _, ok, err := v.modelMakeID(ctx, modelID); err != nil {
		return models.Vehicle{}, err
	} else if !ok {
		return models.Vehicle{}, &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Choose a model from the catalog.",
		}
	}
	if v.ids == nil {
		return models.Vehicle{}, fmt.Errorf("vehicle address minting is not configured")
	}
	id, createdAt, err := v.ids.Next()
	if err != nil {
		return models.Vehicle{}, err
	}
	if err := v.repo.Insert(ctx, id, createdAt, normalized, modelID, entered); err != nil {
		return models.Vehicle{}, err
	}
	logger.FromContext(ctx).Info("recorded vehicle",
		zap.String("layer", "logic"),
		zap.String("op", "RecordVehicle"),
		zap.Int64("vehicle_id", id),
		zap.String("vin", normalized),
	)
	return v.GetVehicle(ctx, strconv.FormatInt(id, 10))
}

// CorrectVehicle changes a non-empty subset of model and entered on.
func (v *Vehicles) CorrectVehicle(ctx context.Context, id string, modelID *int32, enteredOn *string) (models.Vehicle, error) {
	if v.repo == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	parsed, err := ParseVehicleID(id)
	if err != nil {
		return models.Vehicle{}, err
	}
	stored, err := v.repo.Get(ctx, parsed)
	if err != nil {
		return models.Vehicle{}, err
	}
	if stored.StockedOutOn != nil {
		return models.Vehicle{}, &apperrors.WithMessage{
			Err:     apperrors.ErrConflict,
			Message: "The correction is rejected because the vehicle is stocked out.",
		}
	}
	var nextModelID *int32
	var nextEntered *time.Time
	if modelID != nil {
		if err := AcceptModelID(*modelID); err != nil {
			return models.Vehicle{}, err
		}
		if _, ok, err := v.modelMakeID(ctx, *modelID); err != nil {
			return models.Vehicle{}, err
		} else if !ok {
			return models.Vehicle{}, &apperrors.WithMessage{
				Err:     apperrors.ErrValidation,
				Message: "Choose a model from the catalog.",
			}
		}
		nextModelID = modelID
	}
	if enteredOn != nil {
		entered, err := ParseEnteredOn(*enteredOn, v.today())
		if err != nil {
			return models.Vehicle{}, err
		}
		nextEntered = &entered
	}
	if nextModelID == nil && nextEntered == nil {
		return models.Vehicle{}, apperrors.ErrValidation
	}
	if err := v.repo.Correct(ctx, parsed, nextModelID, nextEntered); err != nil {
		return models.Vehicle{}, err
	}
	logger.FromContext(ctx).Info("corrected vehicle",
		zap.String("layer", "logic"),
		zap.String("op", "CorrectVehicle"),
		zap.Int64("vehicle_id", parsed),
	)
	return v.GetVehicle(ctx, strconv.FormatInt(parsed, 10))
}

// StockOutVehicle records stocked-out on. Clearing it is a conflict: the departure stands.
func (v *Vehicles) StockOutVehicle(ctx context.Context, id string, stockedOutOn *string) (models.Vehicle, error) {
	if v.repo == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	parsed, err := ParseVehicleID(id)
	if err != nil {
		return models.Vehicle{}, err
	}
	stored, err := v.repo.Get(ctx, parsed)
	if err != nil {
		return models.Vehicle{}, err
	}
	if stockedOutOn == nil || strings.TrimSpace(*stockedOutOn) == "" {
		if stored.StockedOutOn != nil {
			return models.Vehicle{}, &apperrors.WithMessage{
				Err:     apperrors.ErrConflict,
				Message: "The departure stands.",
			}
		}
		return models.Vehicle{}, apperrors.ErrValidation
	}
	on, err := ParseStockedOutOn(*stockedOutOn, stored.EnteredOn, v.today())
	if err != nil {
		return models.Vehicle{}, err
	}
	if err := v.repo.StockOut(ctx, parsed, on); err != nil {
		return models.Vehicle{}, err
	}
	logger.FromContext(ctx).Info("stocked out vehicle",
		zap.String("layer", "logic"),
		zap.String("op", "StockOutVehicle"),
		zap.Int64("vehicle_id", parsed),
		zap.Time("stocked_out_on", on),
	)
	return v.GetVehicle(ctx, strconv.FormatInt(parsed, 10))
}

// SetProposedAction sets, replaces, or clears the proposed action on aging stock.
func (v *Vehicles) SetProposedAction(ctx context.Context, id string, proposedAction *string) (models.Vehicle, error) {
	if v.repo == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	parsed, err := ParseVehicleID(id)
	if err != nil {
		return models.Vehicle{}, err
	}
	stored, err := v.repo.Get(ctx, parsed)
	if err != nil {
		return models.Vehicle{}, err
	}
	age := aging.InventoryAge(stored.EnteredOn, stored.StockedOutOn, v.today())
	if !aging.IsAging(stored.StockedOutOn, age) {
		return models.Vehicle{}, &apperrors.WithMessage{
			Err:     apperrors.ErrConflict,
			Message: "A proposed action is accepted only for aging stock.",
		}
	}
	next, err := parseProposedAction(proposedAction)
	if err != nil {
		return models.Vehicle{}, err
	}
	if err := v.repo.SetProposedAction(ctx, parsed, next); err != nil {
		return models.Vehicle{}, err
	}
	action := ""
	if next != nil {
		action = *next
	}
	logger.FromContext(ctx).Info("set proposed action",
		zap.String("layer", "logic"),
		zap.String("op", "SetProposedAction"),
		zap.Int64("vehicle_id", parsed),
		zap.String("proposed_action", action),
	)
	return v.GetVehicle(ctx, strconv.FormatInt(parsed, 10))
}

func parseProposedAction(raw *string) (*string, error) {
	if raw == nil || strings.TrimSpace(*raw) == "" {
		return nil, nil
	}
	value := strings.TrimSpace(*raw)
	if _, ok := constant.KnownProposedActions[value]; !ok {
		return nil, &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Choose price reduction planned, wholesale planned, or hold.",
		}
	}
	return &value, nil
}

// GetVehicle returns one vehicle by its address.
func (v *Vehicles) GetVehicle(ctx context.Context, id string) (models.Vehicle, error) {
	if v.repo == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	parsed, err := ParseVehicleID(id)
	if err != nil {
		return models.Vehicle{}, err
	}
	stored, err := v.repo.Get(ctx, parsed)
	if err != nil {
		return models.Vehicle{}, err
	}
	return v.present(stored), nil
}

// ListVehicles returns one page of the inventory.
func (v *Vehicles) ListVehicles(ctx context.Context, makeID, modelID, stockState, stockedOutOn, agingRaw, pageRaw, sizeRaw string) (models.VehiclePage, error) {
	if v.repo == nil {
		return models.VehiclePage{}, fmt.Errorf("vehicles are not configured")
	}
	query, err := ParseList(makeID, modelID, stockState, stockedOutOn, agingRaw, pageRaw, sizeRaw)
	if err != nil {
		return models.VehiclePage{}, err
	}
	if query.MakeID != 0 && query.ModelID != 0 {
		owner, ok, err := v.modelMakeID(ctx, query.ModelID)
		if err != nil {
			return models.VehiclePage{}, err
		}
		if !ok || owner != query.MakeID {
			return models.VehiclePage{}, apperrors.ErrValidation
		}
	}
	bound := aging.OldestInStockDay(v.today()).Format("2006-01-02")
	stored, total, err := v.repo.List(ctx, repositories.ListFilter{
		MakeID:       query.MakeID,
		ModelID:      query.ModelID,
		StockState:   query.StockState,
		StockedOutOn: query.StockedOutOn,
		Aging:        query.Aging,
		AgingBound:   bound,
		Limit:        query.PageSize,
		Offset:       (query.Page - 1) * query.PageSize,
	})
	if err != nil {
		return models.VehiclePage{}, err
	}
	page := models.VehiclePage{
		Vehicles: make([]models.Vehicle, 0, len(stored)),
		Page:     query.Page,
		PageSize: query.PageSize,
		Total:    total,
	}
	for _, item := range stored {
		page.Vehicles = append(page.Vehicles, v.present(item))
	}
	return page, nil
}

func (v *Vehicles) modelMakeID(ctx context.Context, modelID int32) (int32, bool, error) {
	if v.catalog == nil {
		return 0, false, fmt.Errorf("catalog is not configured")
	}
	return v.catalog.ModelMakeID(ctx, modelID)
}

func (v *Vehicles) present(stored models.StoredVehicle) models.Vehicle {
	age := aging.InventoryAge(stored.EnteredOn, stored.StockedOutOn, v.today())
	var stockedOut *string
	if stored.StockedOutOn != nil {
		value := aging.UTCDate(*stored.StockedOutOn).Format("2006-01-02")
		stockedOut = &value
	}
	return models.Vehicle{
		ID:             strconv.FormatInt(stored.ID, 10),
		VIN:            stored.VIN,
		Make:           stored.Make,
		Model:          stored.Model,
		EnteredOn:      aging.UTCDate(stored.EnteredOn).Format("2006-01-02"),
		StockedOutOn:   stockedOut,
		StockState:     aging.StockState(stored.StockedOutOn),
		InventoryAge:   age,
		IsAging:        aging.IsAging(stored.StockedOutOn, age),
		ProposedAction: stored.ProposedAction,
	}
}
