package vehicle

import (
	"strconv"
	"strings"

	"github.com/google/uuid"

	apperrors "inventory_dashboard/internal/common/errors"
)

// ListQuery is a parsed inventory list request.
type ListQuery struct {
	MakeID   string
	ModelID  string
	Page     int
	PageSize int
}

// ParseList applies paging defaults. Page size above 100 is capped at 100.
func ParseList(makeID, modelID, pageRaw, sizeRaw string) (ListQuery, error) {
	query := ListQuery{Page: 1, PageSize: 20}
	if strings.TrimSpace(makeID) != "" {
		if _, err := uuid.Parse(makeID); err != nil {
			return ListQuery{}, apperrors.ErrValidation
		}
		query.MakeID = makeID
	}
	if strings.TrimSpace(modelID) != "" {
		if _, err := uuid.Parse(modelID); err != nil {
			return ListQuery{}, apperrors.ErrValidation
		}
		query.ModelID = modelID
	}
	if pageRaw != "" {
		page, err := strconv.Atoi(pageRaw)
		if err != nil || page < 1 {
			return ListQuery{}, apperrors.ErrValidation
		}
		query.Page = page
	}
	if sizeRaw != "" {
		size, err := strconv.Atoi(sizeRaw)
		if err != nil || size < 1 {
			return ListQuery{}, apperrors.ErrValidation
		}
		if size > 100 {
			size = 100
		}
		query.PageSize = size
	}
	return query, nil
}

// ParseVehicleID requires the vehicle UUID from the address.
func ParseVehicleID(raw string) (string, error) {
	if _, err := uuid.Parse(raw); err != nil {
		return "", apperrors.ErrValidation
	}
	return raw, nil
}
