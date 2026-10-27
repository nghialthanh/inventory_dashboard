package address

import (
	"math"
	"strconv"
	"strings"

	apperrors "inventory_dashboard/internal/common/errors"
)

// Largest integers the API accepts. Catalog addresses fit in a 32-bit integer.
// Vehicle addresses are full positive int64 snowflakes (JSON strings on the wire).
const (
	MaxCatalogID int64 = 2147483647
	MaxVehicleID int64 = math.MaxInt64
)

// Parse reads a positive decimal address up to max. Empty text, a sign, a
// fraction, and a leading zero are rejected.
func Parse(raw string, max int64) (int64, error) {
	text := strings.TrimSpace(raw)
	if text == "" || text[0] < '1' || text[0] > '9' {
		return 0, apperrors.ErrValidation
	}
	for _, ch := range text[1:] {
		if ch < '0' || ch > '9' {
			return 0, apperrors.ErrValidation
		}
	}
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil || value < 1 || value > max {
		return 0, apperrors.ErrValidation
	}
	return value, nil
}
