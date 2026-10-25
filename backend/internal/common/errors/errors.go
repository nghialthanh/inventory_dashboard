package errors

import "errors"

// Sentinel failures the HTTP shell maps onto the shared envelope.
var (
	ErrValidation = errors.New("validation")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)
