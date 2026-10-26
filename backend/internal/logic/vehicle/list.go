package vehicle

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	apperrors "inventory_dashboard/internal/common/errors"
)

// ListQuery is a parsed inventory list request.
type ListQuery struct {
	MakeID       string
	ModelID      string
	StockState   string
	StockedOutOn string
	AgingOnly    bool
	Page         int
	PageSize     int
}

// ParseList applies paging defaults. Page size above 100 is capped at 100.
// aging is omitted to leave the filter off. The only accepted value is true.
// stock_state is omitted to leave the filter off; only in_stock or stocked_out are accepted.
// stocked_out_on is omitted to leave the filter off; when set it must be YYYY-MM-DD.
func ParseList(makeID, modelID, stockState, stockedOutOn, aging, pageRaw, sizeRaw string) (ListQuery, error) {
	query := ListQuery{Page: 1, PageSize: 20}
	if trimmed := strings.TrimSpace(aging); trimmed != "" {
		if trimmed != "true" {
			return ListQuery{}, apperrors.ErrValidation
		}
		query.AgingOnly = true
	}
	if trimmed := strings.TrimSpace(stockState); trimmed != "" {
		if trimmed != "in_stock" && trimmed != "stocked_out" {
			return ListQuery{}, apperrors.ErrValidation
		}
		query.StockState = trimmed
	}
	if trimmed := strings.TrimSpace(stockedOutOn); trimmed != "" {
		if _, err := time.Parse("2006-01-02", trimmed); err != nil {
			return ListQuery{}, apperrors.ErrValidation
		}
		query.StockedOutOn = trimmed
	}
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
