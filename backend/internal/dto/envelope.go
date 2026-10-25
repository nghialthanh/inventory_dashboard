package dto

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Envelope is the shared JSON body for /api/v1 and for /health.
type Envelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Meta carries list paging. Callers omit it on commands that are not lists.
type Meta struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

// Success writes the shared envelope. A nil data value is sent as an empty object.
func Success(c echo.Context, status int, data any) error {
	if data == nil {
		data = map[string]any{}
	}
	return c.JSON(status, Envelope{
		Code:    status,
		Message: "ok",
		Data:    data,
	})
}

// SuccessList writes the shared envelope with paging meta. A nil list is sent as an empty array.
func SuccessList(c echo.Context, data any, meta Meta) error {
	if data == nil {
		data = []any{}
	}
	return c.JSON(http.StatusOK, Envelope{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
		Meta:    &meta,
	})
}

// Error writes the shared envelope with an error_code in data.
func Error(c echo.Context, status int, errorCode, message string) error {
	if message == "" {
		message = "error"
	}
	return c.JSON(status, Envelope{
		Code:    status,
		Message: message,
		Data:    map[string]any{"error_code": errorCode},
	})
}
