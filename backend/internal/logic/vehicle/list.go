package vehicle

import (
	"strconv"
	"strings"
	"time"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/logic/address"
)

// ListQuery is a parsed inventory list request.
type ListQuery struct {
	MakeID       int32
	ModelID      int32
	StockState   string
	StockedOutOn string
	// Aging is nil when omitted, true for aging stock only, false for not aging stock.
	Aging    *bool
	Page     int
	PageSize int
}

// ParseList applies paging defaults. Page size above 100 is capped at 100.
// aging is omitted to leave the filter off; true is aging stock only; false is not aging stock.
// stock_state is omitted to leave the filter off; only in_stock or stocked_out are accepted.
// stocked_out_on is omitted to leave the filter off; when set it must be YYYY-MM-DD.
func ParseList(makeID, modelID, stockState, stockedOutOn, aging, pageRaw, sizeRaw string) (ListQuery, error) {
	query := ListQuery{Page: 1, PageSize: 20}
	if trimmed := strings.TrimSpace(aging); trimmed != "" {
		switch trimmed {
		case "true":
			v := true
			query.Aging = &v
		case "false":
			v := false
			query.Aging = &v
		default:
			return ListQuery{}, apperrors.ErrValidation
		}
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
		parsed, err := address.Parse(makeID, address.MaxCatalogID)
		if err != nil {
			return ListQuery{}, err
		}
		query.MakeID = int32(parsed)
	}
	if strings.TrimSpace(modelID) != "" {
		parsed, err := address.Parse(modelID, address.MaxCatalogID)
		if err != nil {
			return ListQuery{}, err
		}
		query.ModelID = int32(parsed)
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

// ParseVehicleID requires the vehicle address from the URL.
func ParseVehicleID(raw string) (int64, error) {
	return address.Parse(raw, address.MaxVehicleID)
}
