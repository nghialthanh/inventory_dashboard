package middlewares

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"inventory_dashboard/internal/common/trace"
)

// Trace ensures every request has a trace_id on the context and response.
// Accepts an incoming X-Trace-ID or generates a UUID.
func Trace() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get(trace.Header)
			if id == "" {
				id = uuid.NewString()
			}
			ctx := trace.WithID(c.Request().Context(), id)
			c.SetRequest(c.Request().WithContext(ctx))
			c.Response().Header().Set(trace.Header, id)
			return next(c)
		}
	}
}
