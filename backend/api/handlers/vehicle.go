package handlers

import (
	"encoding/json"
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

type correctVehicleRequest struct {
	VIN       *string `json:"vin"`
	ModelID   *string `json:"model_id"`
	EnteredOn *string `json:"entered_on"`
}

func bindCorrection(c echo.Context) (correctVehicleRequest, error) {
	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()
	var req correctVehicleRequest
	if err := decoder.Decode(&req); err != nil {
		return correctVehicleRequest{}, apperrors.ErrValidation
	}
	return req, nil
}

// CorrectVehicle is PATCH /api/v1/vehicles/:id.
func CorrectVehicle(vehicles services.Vehicles) echo.HandlerFunc {
	return func(c echo.Context) error {
		if vehicles == nil {
			return mapError(c, errNotConfigured("vehicles"))
		}
		req, err := bindCorrection(c)
		if err != nil {
			return mapError(c, err)
		}
		vehicle, err := vehicles.CorrectVehicle(c.Request().Context(), c.Param("id"), req.VIN, req.ModelID, req.EnteredOn)
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, vehicle)
	}
}

type stockOutRequest struct {
	StockedOutOn *string `json:"stocked_out_on"`
}

// StockOutVehicle is POST /api/v1/vehicles/:id/stock-out.
func StockOutVehicle(vehicles services.Vehicles) echo.HandlerFunc {
	return func(c echo.Context) error {
		if vehicles == nil {
			return mapError(c, errNotConfigured("vehicles"))
		}
		var req stockOutRequest
		if err := c.Bind(&req); err != nil {
			return mapError(c, apperrors.ErrValidation)
		}
		vehicle, err := vehicles.StockOutVehicle(c.Request().Context(), c.Param("id"), req.StockedOutOn)
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, vehicle)
	}
}

type proposedActionRequest struct {
	ProposedAction *string `json:"proposed_action"`
}

// SetProposedAction is PUT /api/v1/vehicles/:id/proposed-action.
func SetProposedAction(vehicles services.Vehicles) echo.HandlerFunc {
	return func(c echo.Context) error {
		if vehicles == nil {
			return mapError(c, errNotConfigured("vehicles"))
		}
		var req proposedActionRequest
		if err := c.Bind(&req); err != nil {
			return mapError(c, apperrors.ErrValidation)
		}
		vehicle, err := vehicles.SetProposedAction(c.Request().Context(), c.Param("id"), req.ProposedAction)
		if err != nil {
			return mapError(c, err)
		}
		return dto.Success(c, http.StatusOK, vehicle)
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
			c.QueryParam("stock_state"),
			c.QueryParam("stocked_out_on"),
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
