package routes

import (
	"trips-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func TripRoutes(r *gin.Engine, h *handler.TripHandler) {
	trips := r.Group("/api/v1/trips")
	{
		trips.POST("", h.CreateTrip)
		trips.GET("", h.GetTrips)
		trips.GET("/:id", h.GetTripByID)
		trips.PATCH("/:id/assign-driver", h.AssignDriver)
		trips.PATCH("/:id/assign-auto", h.AssignAuto)
		trips.PATCH("/:id/status", h.ChangeTripStatus)
	}
}
