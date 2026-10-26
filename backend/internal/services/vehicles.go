package services

import (
	"context"

	"inventory_dashboard/internal/models"
)

// Vehicles records and reads inventory.
type Vehicles interface {
	RecordVehicle(ctx context.Context, vin, modelID, enteredOn string) (models.Vehicle, error)
	CorrectVehicle(ctx context.Context, id string, vin, modelID, enteredOn *string) (models.Vehicle, error)
	StockOutVehicle(ctx context.Context, id string, stockedOutOn *string) (models.Vehicle, error)
	GetVehicle(ctx context.Context, id string) (models.Vehicle, error)
	ListVehicles(ctx context.Context, makeID, modelID, aging, page, pageSize string) (models.VehiclePage, error)
}
