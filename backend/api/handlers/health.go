package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"inventory_dashboard/internal/dto"
	"inventory_dashboard/internal/services"
)

// Health is GET /health, outside /api/v1.
func Health(health services.Health) echo.HandlerFunc {
	return func(c echo.Context) error {
		if health == nil {
			return dto.Success(c, http.StatusOK, map[string]any{})
		}
		if err := health.Health(); err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, map[string]any{})
	}
}
