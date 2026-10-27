package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"inventory_dashboard/pkg/logger"
)

// AccessLog writes one zap info line per request (skips GET /health).
func AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			if c.Request().Method == "GET" && (c.Path() == "/health" || c.Path() == "/metrics") {
				return err
			}
			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok && he.Code > 0 {
					status = he.Code
				} else if status == 0 {
					status = 500
				}
			}
			logger.FromContext(c.Request().Context()).Info("request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", status),
				zap.Duration("latency", time.Since(start)),
			)
			return err
		}
	}
}
