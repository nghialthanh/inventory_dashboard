package catalog

import (
	"inventory_dashboard/internal/logic/address"
)

// ParseMakeID accepts an empty filter or one make address.
func ParseMakeID(raw string) (*int32, error) {
	if raw == "" {
		return nil, nil
	}
	value, err := address.Parse(raw, address.MaxCatalogID)
	if err != nil {
		return nil, err
	}
	id := int32(value)
	return &id, nil
}
