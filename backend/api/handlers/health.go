package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"inventory_dashboard/internal/dto"
	"inventory_dashboard/internal/services"
)

// Health is GET /health, outside /api/v1.
func Health(svc *services.Services) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := svc.Health(); err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, map[string]any{})
	}
}
