package middlewares

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"inventory_dashboard/pkg/metrics"
)

// Metrics records Prometheus HTTP counters and latency histograms.
func Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			if path == "/metrics" || path == "/health" {
				return next(c)
			}
			start := time.Now()
			err := next(c)
			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok && he.Code > 0 {
					status = he.Code
				} else if status < 400 {
					status = 500
				}
			}
			if path == "" {
				path = c.Request().URL.Path
			}
			method := c.Request().Method
			metrics.HTTPRequests.WithLabelValues(method, path, statusClass(status)).Inc()
			metrics.HTTPDuration.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
			return err
		}
	}
}

func statusClass(status int) string {
	if status < 100 || status > 599 {
		return "unknown"
	}
	return fmt.Sprintf("%dxx", status/100)
}
