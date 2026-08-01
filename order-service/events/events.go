package events

import "time"

const (
	// Имя Exchange для всех бизнес-событий
	ExchangeName = "app.events"

	// Routing keys для заказов
	OrderCreated  = "order.created"
	OrderUpdated  = "order.updated"
	OrderCanceled = "order.canceled"
)

// OrderPayload — структура данных события заказа
type OrderPayload struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	UserEmail string    `json:"user_email,omitempty"`
	CargoType string    `json:"cargo_type,omitempty"`
	Weight    float64   `json:"weight,omitempty"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"` // "created", "in_progress", "completed", "canceled"
	Timestamp time.Time `json:"timestamp"`
}
