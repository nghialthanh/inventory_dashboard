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
	g.PATCH("/vehicles/:id", handlers.CorrectVehicle(vehicles))
	g.POST("/vehicles/:id/stock-out", handlers.StockOutVehicle(vehicles))
	g.PUT("/vehicles/:id/proposed-action", handlers.SetProposedAction(vehicles))
	g.PUT("/vehicles/:id/primary-image", handlers.PutPrimaryImage(vehicles))
	g.DELETE("/vehicles/:id/primary-image", handlers.DeletePrimaryImage(vehicles))
	g.GET("/vehicles/:id/primary-image", handlers.GetPrimaryImage(vehicles))
	g.GET("/vehicles/:id", handlers.GetVehicle(vehicles))
}
