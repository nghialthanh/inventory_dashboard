package errors

import "errors"

// Sentinel failures the HTTP shell maps onto the shared envelope.
var (
	ErrValidation = errors.New("validation")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
)

// WithMessage keeps a sentinel and the text the manager should see.
type WithMessage struct {
	Err     error
	Message string
}

func (e *WithMessage) Error() string { return e.Message }

func (e *WithMessage) Unwrap() error { return e.Err }
