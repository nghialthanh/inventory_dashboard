package vehicle

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"inventory_dashboard/internal/common/constant"
	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/logic/aging"
	"inventory_dashboard/internal/models"
	"inventory_dashboard/internal/repositories"
)

// Vehicles implements the vehicle service.
type Vehicles struct {
	repo    *repositories.Vehicles
	catalog *repositories.Catalog
	today   func() time.Time
}

// New binds vehicle rules to the vehicle and catalog repositories.
// A nil today function uses the process clock.
func New(repo *repositories.Vehicles, catalog *repositories.Catalog, today func() time.Time) *Vehicles {
	if today == nil {
		today = time.Now
	}
	return &Vehicles{repo: repo, catalog: catalog, today: today}
}

// RecordVehicle stores an in-stock vehicle and returns it with derived fields.
func (v *Vehicles) RecordVehicle(ctx context.Context, vin, modelID, enteredOn string) (models.Vehicle, error) {
	if v.repo == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	normalized, err := NormalizeVIN(vin)
	if err != nil {
		return models.Vehicle{}, err
	}
	modelID, err = ParseModelID(modelID)
	if err != nil {
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
	id := uuid.NewString()
	if err := v.repo.Insert(ctx, id, normalized, modelID, entered); err != nil {
		return models.Vehicle{}, err
	}
	return v.GetVehicle(ctx, id)
}

// CorrectVehicle changes a non-empty subset of VIN, model, and entered on.
func (v *Vehicles) CorrectVehicle(ctx context.Context, id string, vin, modelID, enteredOn *string) (models.Vehicle, error) {
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
	var nextVIN, nextModelID *string
	var nextEntered *time.Time
	if vin != nil {
		normalized, err := NormalizeVIN(*vin)
		if err != nil {
			return models.Vehicle{}, err
		}
		nextVIN = &normalized
	}
	if modelID != nil {
		parsedModel, err := ParseModelID(*modelID)
		if err != nil {
			return models.Vehicle{}, err
		}
		if _, ok, err := v.modelMakeID(ctx, parsedModel); err != nil {
			return models.Vehicle{}, err
		} else if !ok {
			return models.Vehicle{}, &apperrors.WithMessage{
				Err:     apperrors.ErrValidation,
				Message: "Choose a model from the catalog.",
			}
		}
		nextModelID = &parsedModel
	}
	if enteredOn != nil {
		entered, err := ParseEnteredOn(*enteredOn, v.today())
		if err != nil {
			return models.Vehicle{}, err
		}
		nextEntered = &entered
	}
	if nextVIN == nil && nextModelID == nil && nextEntered == nil {
		return models.Vehicle{}, apperrors.ErrValidation
	}
	if err := v.repo.Correct(ctx, parsed, nextVIN, nextModelID, nextEntered); err != nil {
		return models.Vehicle{}, err
	}
	return v.GetVehicle(ctx, parsed)
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
	return v.GetVehicle(ctx, parsed)
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
	return v.GetVehicle(ctx, parsed)
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
	if query.MakeID != "" && query.ModelID != "" {
		owner, ok, err := v.modelMakeID(ctx, query.ModelID)
		if err != nil {
			return models.VehiclePage{}, err
		}
		if !ok || owner != query.MakeID {
			return models.VehiclePage{}, apperrors.ErrValidation
		}
	}
	stored, total, err := v.repo.List(ctx, repositories.ListFilter{
		MakeID:        query.MakeID,
		ModelID:       query.ModelID,
		StockState:    query.StockState,
		StockedOutOn:  query.StockedOutOn,
		AgingOnly:     query.AgingOnly,
		EnteredBefore: aging.OldestInStockDay(v.today()).Format("2006-01-02"),
		Limit:         query.PageSize,
		Offset:        (query.Page - 1) * query.PageSize,
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

func (v *Vehicles) modelMakeID(ctx context.Context, modelID string) (string, bool, error) {
	if v.catalog == nil {
		return "", false, fmt.Errorf("catalog is not configured")
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
		ID:              stored.ID,
		VIN:             stored.VIN,
		Make:            stored.Make,
		Model:           stored.Model,
		EnteredOn:       aging.UTCDate(stored.EnteredOn).Format("2006-01-02"),
		StockedOutOn:    stockedOut,
		StockState:      aging.StockState(stored.StockedOutOn),
		InventoryAge:    age,
		IsAging:         aging.IsAging(stored.StockedOutOn, age),
		ProposedAction:  stored.ProposedAction,
		HasPrimaryImage: stored.HasPrimaryImage,
	}
}
