package dto

import "time"

type UpdateOrderRequest struct {
	Status          *int       `json:"status"`
	ManagerID       *int       `json:"manager_id"`
	Price           *float64   `json:"price"`
	PickupAddress   *string    `json:"pickup_address"`
	DeliveryAddress *string    `json:"delivery_address"`
	DateStart       *time.Time `json:"date_start"`
	DateEnd         *time.Time `json:"date_end"`
	PeriodFrom      *time.Time `json:"period_from"`
	PeriodTo        *time.Time `json:"period_to"`
}

type CreateOrderRequest struct {
	// Обязательные поля (должны быть переданы)
	Capacity        int `json:"capacity" binding:"required,min=1"`         // Грузоподъёмность (обязательно, > 0)
	LiftingCapacity int `json:"lifting_capacity" binding:"required,min=1"` // Грузоподъёмность (обязательно, > 0)

	// Опциональные поля (могут быть опущены, тогда в БД будет NULL или значение по умолчанию)
	Type            *int       `json:"type,omitempty"`             // Тип заказа (например, 1, 2, 3)
	ManagerID       *int       `json:"manager_id,omitempty"`       // ID менеджера (может быть назначен позже)
	Status          *int       `json:"status,omitempty"`           // Статус (по умолчанию 1, если не указан)
	Price           *float64   `json:"price,omitempty"`            // Цена (может быть NULL)
	PickupAddress   *string    `json:"pickup_address,omitempty"`   // Адрес забора
	DeliveryAddress *string    `json:"delivery_address,omitempty"` // Адрес доставки
	DateStart       *time.Time `json:"date_start,omitempty"`       // Дата начала
	DateEnd         *time.Time `json:"date_end,omitempty"`         // Дата окончания
	PeriodFrom      *time.Time `json:"period_from,omitempty"`      // Начало периода (если есть)
	PeriodTo        *time.Time `json:"period_to,omitempty"`        // Конец периода
}
