package models

// Make is a manufacturer in the vehicle catalog.
type Make struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Model is a vehicle line that belongs to one make.
type Model struct {
	ID     string `json:"id"`
	MakeID string `json:"make_id"`
	Name   string `json:"name"`
}
