package vehicle

import (
	"regexp"
	"strings"
	"time"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/logic/aging"

	"github.com/google/uuid"
)

var vinPattern = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)

// NormalizeVIN trims and uppercases a VIN, then checks its shape.
func NormalizeVIN(raw string) (string, error) {
	vin := strings.ToUpper(strings.TrimSpace(raw))
	if !vinPattern.MatchString(vin) {
		return "", &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Enter a 17-character VIN using digits and letters except I, O, and Q.",
		}
	}
	return vin, nil
}

// ParseEnteredOn accepts a calendar day on or before today.
func ParseEnteredOn(raw string, today time.Time) (time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}, apperrors.ErrValidation
	}
	entered, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, apperrors.ErrValidation
	}
	entered = aging.UTCDate(entered)
	if entered.After(aging.UTCDate(today)) {
		return time.Time{}, &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Entered on cannot be after today.",
		}
	}
	return entered, nil
}

// ParseStockedOutOn accepts a calendar day on or after entered on and on or before today.
func ParseStockedOutOn(raw string, enteredOn, today time.Time) (time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}, apperrors.ErrValidation
	}
	stockedOutOn, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, apperrors.ErrValidation
	}
	stockedOutOn = aging.UTCDate(stockedOutOn)
	entered := aging.UTCDate(enteredOn)
	if stockedOutOn.Before(entered) {
		return time.Time{}, &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Stocked-out on cannot be before entered on.",
		}
	}
	if stockedOutOn.After(aging.UTCDate(today)) {
		return time.Time{}, &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Stocked-out on cannot be after today.",
		}
	}
	return stockedOutOn, nil
}

// ParseModelID requires a model UUID.
func ParseModelID(raw string) (string, error) {
	if _, err := uuid.Parse(strings.TrimSpace(raw)); err != nil {
		return "", apperrors.ErrValidation
	}
	return strings.TrimSpace(raw), nil
}
