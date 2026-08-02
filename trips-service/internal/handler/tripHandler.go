package handler

import (
	"net/http"
	"strconv"

	"trips-service/dto"
	"trips-service/internal/repository/model"
	"trips-service/internal/service"

	"github.com/gin-gonic/gin"
)

type TripHandler struct {
	service service.TripService
}

func NewTripHandler(service service.TripService) *TripHandler {
	return &TripHandler{service: service}
}

func (h *TripHandler) CreateTrip(c *gin.Context) {
	var input dto.CreateTripInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trip := &model.Trip{
		OrderID: input.OrderID,
		Status:  "ASSIGNED",
	}

	if err := h.service.CreateTrip(c.Request.Context(), trip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trip)
}

// GetTripByID GET /api/v1/trips/:id
func (h *TripHandler) GetTripByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	trip, err := h.service.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

// GetTrips GET /api/v1/trips?page=1&page_size=20&driver_id=5&auto_id=12&status=IN_TRANSIT
func (h *TripHandler) GetTrips(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := &dto.TripFilter{}

	if driverIDStr := c.Query("driver_id"); driverIDStr != "" {
		filter.DriverID = &driverIDStr
	}

	if autoIDStr := c.Query("auto_id"); autoIDStr != "" {
		filter.AutoID = &autoIDStr
	}

	if status := c.Query("status"); status != "" {
		filter.Statuses = []string{status}
	}

	trips, total, err := h.service.GetTrips(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  trips,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// AssignDriver PATCH /api/v1/trips/:id/assign-driver
func (h *TripHandler) AssignDriver(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	var input dto.UpdateTripInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignDriver(c.Request.Context(), id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "driver assigned successfully"})
}

// AssignAuto PATCH /api/v1/trips/:id/assign-auto
func (h *TripHandler) AssignAuto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	var input dto.UpdateTripInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignAuto(c.Request.Context(), id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "auto assigned successfully"})
}

type ChangeStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ChangeTripStatus PATCH /api/v1/trips/:id/status
func (h *TripHandler) ChangeTripStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip id"})
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangeTripStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}
