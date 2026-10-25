package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	apperrors "inventory_dashboard/internal/common/errors"
	"inventory_dashboard/internal/dto"
	"inventory_dashboard/internal/services"
)

type recordVehicleRequest struct {
	VIN       string `json:"vin"`
	ModelID   string `json:"model_id"`
	EnteredOn string `json:"entered_on"`
}

// RecordVehicle is POST /api/v1/vehicles.
func RecordVehicle(vehicles services.Vehicles) echo.HandlerFunc {
	return func(c echo.Context) error {
		if vehicles == nil {
			return mapError(c, errNotConfigured("vehicles"))
		}
		var req recordVehicleRequest
		if err := c.Bind(&req); err != nil {
			return mapError(c, apperrors.ErrValidation)
		}
		vehicle, err := vehicles.RecordVehicle(c.Request().Context(), req.VIN, req.ModelID, req.EnteredOn)
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusCreated, vehicle)
	}
}

// GetVehicle is GET /api/v1/vehicles/:id.
func GetVehicle(vehicles services.Vehicles) echo.HandlerFunc {
	return func(c echo.Context) error {
		if vehicles == nil {
			return mapError(c, errNotConfigured("vehicles"))
		}
		vehicle, err := vehicles.GetVehicle(c.Request().Context(), c.Param("id"))
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, vehicle)
	}
}

// ListVehicles is GET /api/v1/vehicles.
func ListVehicles(vehicles services.Vehicles) echo.HandlerFunc {
	return func(c echo.Context) error {
		if vehicles == nil {
			return mapError(c, errNotConfigured("vehicles"))
		}
		page, err := vehicles.ListVehicles(
			c.Request().Context(),
			c.QueryParam("make_id"),
			c.QueryParam("model_id"),
			c.QueryParam("aging"),
			c.QueryParam("page"),
			c.QueryParam("page_size"),
		)
		if err != nil {
			return mapError(c, err)
		}
		return dto.SuccessList(c, page.Vehicles, dto.Meta{
			Page:     page.Page,
			PageSize: page.PageSize,
			Total:    page.Total,
		})
	}
}
