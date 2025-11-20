package routes

import (
	"yard-planning/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *handler.Handler) {

	// API Grouping (optional)
	api := e.Group("/api")

	// Yard Planning Endpoints
	api.POST("/suggestion", h.Suggest)
	api.POST("/placement", h.Placement)
	api.POST("/pickup", h.Pickup)
}
