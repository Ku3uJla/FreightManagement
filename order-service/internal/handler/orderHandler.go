// internal/handler/order_handler.go
package handler

import (
	"net/http"
	"strconv"

	"order-service/internal/dto"
	"order-service/internal/repository/model"
	"order-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

const DefaultPageSize = 20

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req dto.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.orderService.UpdateOrder(ctx, id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var status int
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.orderService.UpdateStatus(c.Request.Context(), id, status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *OrderHandler) AssignManager(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	managerID := c.MustGet("userID").(int)

	if err := h.orderService.AssignManager(c.Request.Context(), id, managerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "manager assigned"})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	var filter dto.OrderFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter: " + err.Error()})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	orders, total, err := h.orderService.GetOrders(c.Request.Context(), &filter, page, DefaultPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if orders == nil {
		orders = []model.Order{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":       page,
			"pageSize":   DefaultPageSize,
			"total":      total,
			"totalPages": (total + int64(DefaultPageSize) - 1) / int64(DefaultPageSize),
		},
	})
}

func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	orders, total, err := h.orderService.GetOrdersByUser(c.Request.Context(), userID, page, DefaultPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orders == nil {
		orders = []model.Order{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":       page,
			"pageSize":   DefaultPageSize,
			"total":      total,
			"totalPages": (total + int64(DefaultPageSize) - 1) / int64(DefaultPageSize),
		},
	})
}

func (h *OrderHandler) GetOrdersByDriver(c *gin.Context) {
	driverID, err := strconv.Atoi(c.Param("driver_id"))
	if err != nil || driverID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	orders, err := h.orderService.GetOrdersByDriver(c.Request.Context(), driverID, page, DefaultPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if orders == nil {
		orders = []model.Order{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":     page,
			"pageSize": DefaultPageSize,
		},
	})
}
