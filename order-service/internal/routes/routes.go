package routes

import (
	"order-service/internal/handler"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine, h *handler.OrderHandler) {
	r.POST("/orders", h.CreateOrder)
	r.GET("/orders/:id", h.GetOrderByID)
	r.PUT("/orders/:id", h.UpdateOrder)
	r.PATCH("/orders/:id/status", h.UpdateStatus)
	r.POST("/orders/:id/manager", h.AssignManager)
	r.GET("/orders", h.GetOrders)

	r.GET("/users/:user_id/orders", h.GetOrdersByUser)
	r.GET("/drivers/:driver_id/orders", h.GetOrdersByDriver)
}
