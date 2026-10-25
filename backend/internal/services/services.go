package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/logic/aging"
	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/logic/vehicle"
	"inventory_dashboard/internal/models"
	"inventory_dashboard/internal/repositories"
)

// Deps are the stores and clock the process registers.
type Deps struct {
	Catalog  *repositories.Catalog
	Vehicles *repositories.Vehicles
	Today    func() time.Time
}

// Services is the application registry. The process constructs one in Run.
type Services struct {
	catalog  *repositories.Catalog
	vehicles *repositories.Vehicles
	today    func() time.Time
}

// New returns the application services. A nil store leaves that command unavailable.
func New(deps Deps) *Services {
	if deps.Today == nil {
		deps.Today = time.Now
	}
	return &Services{catalog: deps.Catalog, vehicles: deps.Vehicles, today: deps.Today}
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

// RecordVehicle stores an in-stock vehicle and returns it with derived fields.
func (s *Services) RecordVehicle(ctx context.Context, vin, modelID, enteredOn string) (models.Vehicle, error) {
	if s.vehicles == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	normalized, err := vehicle.NormalizeVIN(vin)
	if err != nil {
		return models.Vehicle{}, err
	}
	modelID, err = vehicle.ParseModelID(modelID)
	if err != nil {
		return models.Vehicle{}, err
	}
	entered, err := vehicle.ParseEnteredOn(enteredOn, s.today())
	if err != nil {
		return models.Vehicle{}, err
	}
	if _, ok, err := s.vehicles.ModelMakeID(ctx, modelID); err != nil {
		return models.Vehicle{}, err
	} else if !ok {
		return models.Vehicle{}, &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Choose a model from the catalog.",
		}
	}
	id := uuid.NewString()
	if err := s.vehicles.Insert(ctx, id, normalized, modelID, entered); err != nil {
		return models.Vehicle{}, err
	}
	return s.GetVehicle(ctx, id)
}

// GetVehicle returns one vehicle by its address.
func (s *Services) GetVehicle(ctx context.Context, id string) (models.Vehicle, error) {
	if s.vehicles == nil {
		return models.Vehicle{}, fmt.Errorf("vehicles are not configured")
	}
	parsed, err := vehicle.ParseVehicleID(id)
	if err != nil {
		return models.Vehicle{}, err
	}
	stored, err := s.vehicles.Get(ctx, parsed)
	if err != nil {
		return models.Vehicle{}, err
	}
	return s.present(stored), nil
}

// ListVehicles returns one page of the inventory.
func (s *Services) ListVehicles(ctx context.Context, makeID, modelID, pageRaw, sizeRaw string) (models.VehiclePage, error) {
	if s.vehicles == nil {
		return models.VehiclePage{}, fmt.Errorf("vehicles are not configured")
	}
	query, err := vehicle.ParseList(makeID, modelID, pageRaw, sizeRaw)
	if err != nil {
		return models.VehiclePage{}, err
	}
	if query.MakeID != "" && query.ModelID != "" {
		owner, ok, err := s.vehicles.ModelMakeID(ctx, query.ModelID)
		if err != nil {
			return models.VehiclePage{}, err
		}
		if !ok || owner != query.MakeID {
			return models.VehiclePage{}, apperrors.ErrValidation
		}
	}
	stored, total, err := s.vehicles.List(ctx, repositories.ListFilter{
		MakeID:  query.MakeID,
		ModelID: query.ModelID,
		Limit:   query.PageSize,
		Offset:  (query.Page - 1) * query.PageSize,
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
		page.Vehicles = append(page.Vehicles, s.present(item))
	}
	return page, nil
}

func (s *Services) present(stored models.StoredVehicle) models.Vehicle {
	age := aging.InventoryAge(stored.EnteredOn, stored.StockedOutOn, s.today())
	var stockedOut *string
	if stored.StockedOutOn != nil {
		value := aging.UTCDate(*stored.StockedOutOn).Format("2006-01-02")
		stockedOut = &value
	}
	state := aging.StockState(stored.StockedOutOn)
	return models.Vehicle{
		ID:              stored.ID,
		VIN:             stored.VIN,
		Make:            stored.Make,
		Model:           stored.Model,
		EnteredOn:       aging.UTCDate(stored.EnteredOn).Format("2006-01-02"),
		StockedOutOn:    stockedOut,
		StockState:      state,
		InventoryAge:    age,
		IsAging:         aging.IsAging(stored.StockedOutOn, age),
		ProposedAction:  stored.ProposedAction,
		HasPrimaryImage: stored.HasPrimaryImage,
	}
}
