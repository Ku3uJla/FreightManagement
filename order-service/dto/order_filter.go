package dto

import "time"

type OrderFilter struct {
	// Основные фильтры
	UserID    *int `form:"user_id" json:"user_id"`
	Status    *int `form:"status" json:"status"`
	ManagerID *int `form:"manager_id" json:"manager_id"`
	Type      *int `form:"type" json:"type"`

	// Фильтры по вместимости (диапазон)
	CapacityMin        *int `form:"capacity_min" json:"capacity_min"`
	CapacityMax        *int `form:"capacity_max" json:"capacity_max"`
	LiftingCapacityMin *int `form:"lifting_capacity_min" json:"lifting_capacity_min"`
	LiftingCapacityMax *int `form:"lifting_capacity_max" json:"lifting_capacity_max"`

	// Фильтры по цене (диапазон)
	PriceMin *float64 `form:"price_min" json:"price_min"`
	PriceMax *float64 `form:"price_max" json:"price_max"`

	// Фильтры по датам (диапазоны)
	DateCreateFrom *time.Time `form:"date_create_from" json:"date_create_from"`
	DateCreateTo   *time.Time `form:"date_create_to" json:"date_create_to"`
	DateStartFrom  *time.Time `form:"date_start_from" json:"date_start_from"`
	DateStartTo    *time.Time `form:"date_start_to" json:"date_start_to"`
	DateEndFrom    *time.Time `form:"date_end_from" json:"date_end_from"`
	DateEndTo      *time.Time `form:"date_end_to" json:"date_end_to"`
	PeriodFromFrom *time.Time `form:"period_from_from" json:"period_from_from"`
	PeriodFromTo   *time.Time `form:"period_from_to" json:"period_from_to"`
	PeriodToFrom   *time.Time `form:"period_to_from" json:"period_to_from"`
	PeriodToTo     *time.Time `form:"period_to_to" json:"period_to_to"`

	// Фильтры по адресам (частичное совпадение)
	PickupAddress   *string `form:"pickup_address" json:"pickup_address"`
	DeliveryAddress *string `form:"delivery_address" json:"delivery_address"`

	// Фильтры по связанным таблицам (через driver_auto)
	DriverID *int `form:"driver_id" json:"driver_id"`
	AutoID   *int `form:"auto_id" json:"auto_id"`

	// Пагинация
	Limit  int `form:"limit" json:"limit" default:"20"`
	Offset int `form:"offset" json:"offset" default:"0"`
}
