package events

import "time"

const (
	UserCreated  = "user.created"
	OrderCreated = "order.created"
	OrderUpdated = "order.updated"
	DriverOrder  = "driver.order"
)

type UserCreatedPayload struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderPayload struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	UserEmail string    `json:"user_email"`
	DriverID  string    `json:"driver_id,omitempty"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DriverOrderPayload struct {
	DriverID    string    `json:"driver_id"`
	DriverEmail string    `json:"driver_email"`
	OrderID     string    `json:"order_id"`
	Message     string    `json:"message"`
	AssignedAt  time.Time `json:"assigned_at"`
}
