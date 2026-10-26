package services

import (
	"context"
	"io"

	"inventory_dashboard/internal/models"
)

// Vehicles records and reads inventory.
type Vehicles interface {
	RecordVehicle(ctx context.Context, vin, modelID, enteredOn string) (models.Vehicle, error)
	CorrectVehicle(ctx context.Context, id string, vin, modelID, enteredOn *string) (models.Vehicle, error)
	StockOutVehicle(ctx context.Context, id string, stockedOutOn *string) (models.Vehicle, error)
	SetProposedAction(ctx context.Context, id string, proposedAction *string) (models.Vehicle, error)
	PutPrimaryImage(ctx context.Context, id string, size int64, body io.Reader) (models.Vehicle, error)
	DeletePrimaryImage(ctx context.Context, id string) (models.Vehicle, error)
	OpenPrimaryImage(ctx context.Context, id string) (body io.ReadCloser, contentType string, err error)
	GetVehicle(ctx context.Context, id string) (models.Vehicle, error)
	ListVehicles(ctx context.Context, makeID, modelID, stockState, stockedOutOn, aging, page, pageSize string) (models.VehiclePage, error)
}
