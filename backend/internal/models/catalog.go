package models

// Make is a manufacturer in the vehicle catalog.
type Make struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// Model is a vehicle line that belongs to one make.
type Model struct {
	ID     int32  `json:"id"`
	MakeID int32  `json:"make_id"`
	Name   string `json:"name"`
}
