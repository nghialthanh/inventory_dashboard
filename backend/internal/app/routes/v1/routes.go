package v1

import (
	"github.com/labstack/echo/v4"

	"inventory_dashboard/api/handlers"
	"inventory_dashboard/internal/services"
)

// Register mounts versioned inventory routes on /api/v1.
// Health stays on the root router, outside this group.
func Register(g *echo.Group, catalog services.Catalog, vehicles services.Vehicles) {
	g.GET("/catalog/makes", handlers.ListMakes(catalog))
	g.GET("/catalog/models", handlers.ListModels(catalog))
	g.GET("/vehicles", handlers.ListVehicles(vehicles))
	g.POST("/vehicles", handlers.RecordVehicle(vehicles))
	g.GET("/vehicles/:id", handlers.GetVehicle(vehicles))
}
