package events

import "time"

const (
	ExchangeName = "app.events"
	UserCreated  = "user.created"
)

// UserCreatedPayload — структура события создания пользователя
type UserCreatedPayload struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"` // "client", "driver", "admin"
	CreatedAt time.Time `json:"created_at"`
}
