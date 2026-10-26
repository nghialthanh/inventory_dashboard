package models

import "time"

// Ref is an id and display name embedded on a vehicle.
type Ref struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Vehicle is the list and detail shape.
type Vehicle struct {
	ID             int64   `json:"id"`
	VIN            string  `json:"vin"`
	Make           Ref     `json:"make"`
	Model          Ref     `json:"model"`
	EnteredOn      string  `json:"entered_on"`
	StockedOutOn   *string `json:"stocked_out_on"`
	StockState     string  `json:"stock_state"`
	InventoryAge   int     `json:"inventory_age"`
	IsAging        bool    `json:"is_aging"`
	ProposedAction *string `json:"proposed_action"`
}

// StoredVehicle is a vehicle row before derived fields are applied.
type StoredVehicle struct {
	ID             int64
	VIN            string
	Make           Ref
	Model          Ref
	EnteredOn      time.Time
	StockedOutOn   *time.Time
	ProposedAction *string
}

// VehiclePage is one page of the inventory list.
type VehiclePage struct {
	Vehicles []Vehicle
	Page     int
	PageSize int
	Total    int
}
