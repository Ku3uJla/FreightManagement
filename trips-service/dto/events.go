package dto

import "time"

// Константы для EventType
const (
	EventTypeTripCreated       = "trip.created"
	EventTypeTripStatusUpdated = "trip.status_updated"
	EventTypeTripCompleted     = "trip.completed"
	EventTypeTripCancelled     = "trip.cancelled"
)

// EventEnvelope — универсальная обертка для всех событий в RabbitMQ
type EventEnvelope struct {
	EventID   string      `json:"event_id"`   // Уникальный UUID события (для дедупликации)
	EventType string      `json:"event_type"` // Тип события (например: "trip.status_updated")
	Timestamp time.Time   `json:"timestamp"`  // Время генерации события (UTC)
	Payload   interface{} `json:"payload"`    // Полезная нагрузка
}

// TripStatusUpdatedPayload — данные при смене промежуточного статуса рейса
type TripStatusUpdatedPayload struct {
	TripID    string    `json:"trip_id"`
	OrderID   string    `json:"order_id"`
	DriverID  string    `json:"driver_id"`
	AutoID    string    `json:"auto_id"`
	OldStatus string    `json:"old_status"` // Предыдущий статус (например: "ASSIGNED")
	NewStatus string    `json:"new_status"` // Новый статус (например: "AT_LOADING", "IN_TRANSIT", "AT_UNLOADING")
	UpdatedAt time.Time `json:"updated_at"`
}

// TripCompletedPayload — данные при успешном завершении рейса
type TripCompletedPayload struct {
	TripID     string    `json:"trip_id"`
	OrderID    string    `json:"order_id"`
	DriverID   string    `json:"driver_id"`
	AutoID     string    `json:"auto_id"`
	FinishedAt time.Time `json:"finished_at"`
}

// TripCancelledPayload — данные при отмене рейса
type TripCancelledPayload struct {
	TripID      string    `json:"trip_id"`
	OrderID     string    `json:"order_id"`
	DriverID    string    `json:"driver_id"`
	AutoID      string    `json:"auto_id"`
	Reason      string    `json:"reason"`       // Причина отмены
	CancelledBy string    `json:"cancelled_by"` // "DRIVER", "DISPATCHER", "CLIENT"
	CancelledAt time.Time `json:"cancelled_at"`
}
