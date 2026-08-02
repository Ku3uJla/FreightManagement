package dto

import "time"

// CreateTripInput — DTO для входящего запроса на создание рейса
type CreateTripInput struct {
	OrderID int `json:"order_id" binding:"required"`

	// Поля планирования (опционально, если логисты планируют выезд заранее)
	PlannedStartAt *time.Time `json:"planned_start_at,omitempty"`
}

// UpdateTripInput — DTO для частичного обновления рейса (PATCH)
type UpdateTripInput struct {
	DriverID *string `json:"driver_id"`
	AutoID   *string `json:"auto_id"`
	Comment  *string `json:"comment"`
}
