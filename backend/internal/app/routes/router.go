package routes

import (
	"github.com/labstack/echo/v4"

	"inventory_dashboard/api/handlers"
	"inventory_dashboard/internal/app/routes/v1"
	"inventory_dashboard/internal/services"
)

// Config is the HTTP shell. An empty StaticDir serves the API only.
type Config struct {
	StaticDir string
	Services  *services.Services
}

// New builds the HTTP server. GET /health is registered outside /api/v1.
func New(cfg Config) *echo.Echo {
	if cfg.Services == nil {
		cfg.Services = services.New(services.Deps{})
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/health", handlers.Health(cfg.Services))
	v1.Register(e.Group("/api/v1"), cfg.Services)
	mountStatic(e, cfg.StaticDir)
	return e
}
