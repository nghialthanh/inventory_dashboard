package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/dto"
	"inventory_dashboard/pkg/logger"
)

func mapError(c echo.Context, err error) error {
	message := "error"
	var with *apperrors.WithMessage
	if errors.As(err, &with) && with.Message != "" {
		message = with.Message
	}
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		return dto.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", message)
	case errors.Is(err, apperrors.ErrNotFound):
		return dto.Error(c, http.StatusNotFound, "NOT_FOUND", message)
	case errors.Is(err, apperrors.ErrConflict):
		return dto.Error(c, http.StatusConflict, "CONFLICT", message)
	default:
		logger.FromContext(c.Request().Context()).Error("handler internal error",
			zap.String("layer", "handler"),
			zap.String("path", c.Path()),
			zap.Error(err),
		)
		return dto.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
	}
}

func errNotConfigured(name string) error {
	return fmt.Errorf("%s is not configured", name)
}
