package routes

import (
	"resource-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func DriverRoutes(r *gin.Engine, h *handler.DriverHandler) {
	r.POST("/drivers", h.CreateDriver)
	r.POST("/drivers/categories", h.CreateDriverCategory)
	r.GET("/drivers/:id", h.GetDriverByID)
	r.GET("/drivers/:id/categories", h.GetDriverCategories)
	r.GET("/drivers", h.GetDriversByFilter)
	r.PUT("/drivers/:id/status", h.UpdateDriverStatus)
}

func AutoRoutes(r *gin.Engine, h handler.AutoHandler) {
	r.GET("/autos/:id", h.GetAutoByID)
	r.GET("/autos", h.GetAutosWithFilter)
	r.PUT("/autos/:id/status", h.UpdateAutoStatus)
	r.POST("/autos", h.CreateAuto)
}
