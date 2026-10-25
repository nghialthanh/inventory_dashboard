package catalog

import (
	"github.com/google/uuid"

	apperrors "inventory_dashboard/internal/common/errors"
)

// ParseMakeID accepts an empty filter or one make UUID.
func ParseMakeID(raw string) (*string, error) {
	if raw == "" {
		return nil, nil
	}
	if _, err := uuid.Parse(raw); err != nil {
		return nil, apperrors.ErrValidation
	}
	return &raw, nil
}
