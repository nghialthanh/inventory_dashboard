package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"inventory_dashboard/api/handlers"
	"inventory_dashboard/api/middlewares"
	"inventory_dashboard/internal/app/routes/v1"
	"inventory_dashboard/internal/services"
)

// Config is the HTTP adapter. An empty StaticDir serves the API only.
// Health, Catalog, and Vehicles are service interfaces implemented by logic.
type Config struct {
	StaticDir string
	Health    services.Health
	Catalog   services.Catalog
	Vehicles  services.Vehicles
}

// New builds the HTTP server. GET /health and GET /metrics are registered outside /api/v1.
func New(cfg Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middlewares.Trace(), middlewares.Metrics(), middlewares.AccessLog())

	e.GET("/health", handlers.Health(cfg.Health))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	v1.Register(e.Group("/api/v1"), cfg.Catalog, cfg.Vehicles)
	mountStatic(e, cfg.StaticDir)
	return e
}
