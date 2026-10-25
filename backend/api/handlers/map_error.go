package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/dto"
)

func mapError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		return dto.Error(c, http.StatusBadRequest, "VALIDATION_ERROR")
	case errors.Is(err, apperrors.ErrNotFound):
		return dto.Error(c, http.StatusNotFound, "NOT_FOUND")
	case errors.Is(err, apperrors.ErrConflict):
		return dto.Error(c, http.StatusConflict, "CONFLICT")
	default:
		return dto.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR")
	}
}
