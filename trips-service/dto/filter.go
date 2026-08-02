// internal/repository/model/filter.go
package dto

import "time"

type TripFilter struct {
	DriverID *string    `json:"driver_id"`
	AutoID   *string    `json:"auto_id"`
	Statuses []string   `json:"statuses"`  // Поиск по нескольким статусам
	DateFrom *time.Time `json:"date_from"` // Диапазон дат
	DateTo   *time.Time `json:"date_to"`
}
