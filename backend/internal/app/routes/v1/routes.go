package v1

import (
	"github.com/labstack/echo/v4"

	"inventory_dashboard/api/handlers"
	"inventory_dashboard/internal/services"
)

// Register mounts versioned inventory routes on /api/v1.
// Health stays on the root router, outside this group.
func Register(g *echo.Group, svc *services.Services) {
	g.GET("/catalog/makes", handlers.ListMakes(svc))
	g.GET("/catalog/models", handlers.ListModels(svc))
	g.GET("/vehicles", handlers.ListVehicles(svc))
	g.POST("/vehicles", handlers.RecordVehicle(svc))
	g.GET("/vehicles/:id", handlers.GetVehicle(svc))
}
