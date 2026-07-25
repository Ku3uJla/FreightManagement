package handler

import (
	"net/http"
	"strconv"

	"user-service/internal/filters"
	"user-service/internal/repository/model"

	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AutoHandler struct {
	autoService *service.AutoService
}

func NewAutoHandler(autoService *service.AutoService) *AutoHandler {
	return &AutoHandler{autoService: autoService}
}

// GetAutoByID GET /autos/:id
func (h *AutoHandler) GetAutoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auto id"})
		return
	}

	auto, err := h.autoService.GetAuto(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, auto)
}

// GetAutosWithFilter GET /autos?capacity=...&lifting_capacity=...&status=...
func (h *AutoHandler) GetAutosWithFilter(c *gin.Context) {
	var filter filters.AutoFilter

	if capStr := c.Query("capacity"); capStr != "" {
		capVal, err := strconv.Atoi(capStr)
		if err == nil {
			filter.Capacity = &capVal
		}
	}
	if liftStr := c.Query("lifting_capacity"); liftStr != "" {
		liftVal, err := strconv.Atoi(liftStr)
		if err == nil {
			filter.LiftingCapacity = &liftVal
		}
	}
	if statusStr := c.Query("status"); statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err == nil {
			filter.Status = &statusVal
		}
	}

	autos, err := h.autoService.GetAutosWithFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if autos == nil {
		autos = &[]model.Auto{} // возвращаем пустой массив, а не nil
	}
	c.JSON(http.StatusOK, autos)
}

// UpdateAutoStatus PUT /autos/:id/status
// Тело: { "status": 3 }
func (h *AutoHandler) UpdateAutoStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auto id"})
		return
	}

	var req struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Status < 1 || req.Status > 3 { // пример допустимых статусов
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be 1,2,3"})
		return
	}

	if err := h.autoService.UpdateStatusAuto(c.Request.Context(), req.Status, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// CreateAuto POST /autos
// Тело: поля модели Auto
func (h *AutoHandler) CreateAuto(c *gin.Context) {
	var auto model.Auto
	if err := c.ShouldBindJSON(&auto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.autoService.CreateAuto(c.Request.Context(), &auto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "auto created", "id": auto.ID}) // предполагаем, что ID заполняется после создания
}

// Регистрация маршрутов (можно вынести в отдельную функцию)
//func (h *AutoHandler) RegisterRoutes(r *gin.RouterGroup) {
//	r.GET("/autos/:id", h.GetAutoByID)
//	r.GET("/autos", h.GetAutosWithFilter)
//	r.PUT("/autos/:id/status", h.UpdateAutoStatus)
//	r.POST("/autos", h.CreateAuto)
//}
