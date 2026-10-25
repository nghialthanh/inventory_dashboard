package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"inventory_dashboard/internal/dto"
	"inventory_dashboard/internal/services"
)

// ListMakes is GET /api/v1/catalog/makes.
func ListMakes(catalog services.Catalog) echo.HandlerFunc {
	return func(c echo.Context) error {
		if catalog == nil {
			return mapError(c, errNotConfigured("catalog"))
		}
		makes, err := catalog.ListMakes(c.Request().Context())
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, makes)
	}
}

// ListModels is GET /api/v1/catalog/models.
func ListModels(catalog services.Catalog) echo.HandlerFunc {
	return func(c echo.Context) error {
		if catalog == nil {
			return mapError(c, errNotConfigured("catalog"))
		}
		list, err := catalog.ListModels(c.Request().Context(), c.QueryParam("make_id"))
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, list)
	}
}
