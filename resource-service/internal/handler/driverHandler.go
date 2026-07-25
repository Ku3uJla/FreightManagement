package handler

import (
	"net/http"
	"strconv"

	"resource-service/internal/filters"
	"resource-service/internal/repository/model"
	"resource-service/internal/service"

	"github.com/gin-gonic/gin"
)

type DriverHandler struct {
	driverService *service.DriverService
}

func NewDriverHandler(driverService *service.DriverService) *DriverHandler {
	return &DriverHandler{driverService: driverService}
}

// CreateDriver POST /drivers
// Тело: model.Driver (без ID, он генерируется)
func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var driver model.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.driverService.CreateDriver(c.Request.Context(), &driver); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "driver created", "id": driver.ID})
}

// CreateDriverCategory POST /drivers/categories
// Тело: model.DriverCategory (содержит driver_id и category)
func (h *DriverHandler) CreateDriverCategory(c *gin.Context) {
	var category model.DriverCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.driverService.CreateDriverCategory(c.Request.Context(), &category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "driver category added"})
}

// GetDriverByID GET /drivers/:id
func (h *DriverHandler) GetDriverByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	driver, err := h.driverService.GetDriverByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, driver)
}

// GetDriverCategories GET /drivers/:id/categories
func (h *DriverHandler) GetDriverCategories(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	categories, err := h.driverService.GetDriverCategories(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories) // categories уже []model.DriverCategory (может быть пустым)
}

// GetDriversByFilter GET /drivers?status=...&category=...
func (h *DriverHandler) GetDriversByFilter(c *gin.Context) {
	var filter filters.DriverFilter

	if statusStr := c.Query("status"); statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err == nil {
			filter.Status = &statusVal
		}
	}
	if categoryStr := c.Query("category"); categoryStr != "" {
		filter.Category = &categoryStr
	}

	drivers, err := h.driverService.GetDriversByFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if drivers == nil {
		drivers = []model.Driver{} // пустой массив вместо nil
	}
	c.JSON(http.StatusOK, drivers)
}

// UpdateDriverStatus PUT /drivers/:id/status
// Тело: { "status": 2 }
func (h *DriverHandler) UpdateDriverStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver id"})
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Валидация статуса может быть также в сервисе, но дублируем здесь для раннего ответа
	if req.Status < 1 || req.Status > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 1, 2, or 3"})
		return
	}

	if err := h.driverService.UpdateDriverStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "driver status updated"})
}

// RegisterRoutes регистрирует маршруты для DriverHandler
func (h *DriverHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/drivers", h.CreateDriver)
	r.POST("/drivers/categories", h.CreateDriverCategory)
	r.GET("/drivers/:id", h.GetDriverByID)
	r.GET("/drivers/:id/categories", h.GetDriverCategories)
	r.GET("/drivers", h.GetDriversByFilter)
	r.PUT("/drivers/:id/status", h.UpdateDriverStatus)
}
