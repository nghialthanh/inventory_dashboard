package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"inventory_dashboard/internal/dto"
	"inventory_dashboard/internal/services"
)

// ListMakes is GET /api/v1/catalog/makes.
func ListMakes(svc *services.Services) echo.HandlerFunc {
	return func(c echo.Context) error {
		makes, err := svc.ListMakes(c.Request().Context())
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, makes)
	}
}

// ListModels is GET /api/v1/catalog/models.
func ListModels(svc *services.Services) echo.HandlerFunc {
	return func(c echo.Context) error {
		list, err := svc.ListModels(c.Request().Context(), c.QueryParam("make_id"))
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, list)
	}
}
