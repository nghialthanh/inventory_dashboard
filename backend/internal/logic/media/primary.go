package media

import (
	"fmt"
	"net/http"

	"inventory_dashboard/internal/common/constant"
	apperrors "inventory_dashboard/internal/common/errors"
)

// ObjectKey is the MinIO object key for a vehicle's primary image.
func ObjectKey(vehicleID string) string {
	return fmt.Sprintf("vehicles/%s/primary", vehicleID)
}

// DetectContentType sniffs image bytes and returns a known type or an error.
func DetectContentType(data []byte) (string, error) {
	detected := http.DetectContentType(data)
	switch detected {
	case constant.ContentTypeJPEG, constant.ContentTypePNG:
		return detected, nil
	default:
		if isWebP(data) {
			return constant.ContentTypeWebP, nil
		}
		return "", &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Upload a JPEG, PNG, or WebP image of at most 5 MiB.",
		}
	}
}

func isWebP(header []byte) bool {
	return len(header) >= 12 &&
		string(header[0:4]) == "RIFF" &&
		string(header[8:12]) == "WEBP"
}

// ValidateSize rejects uploads larger than the primary-image limit.
func ValidateSize(size int64) error {
	if size <= 0 || size > constant.PrimaryImageMaxBytes {
		return &apperrors.WithMessage{
			Err:     apperrors.ErrValidation,
			Message: "Upload a JPEG, PNG, or WebP image of at most 5 MiB.",
		}
	}
	return nil
}
