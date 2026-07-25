package model

import (
    "time"
)

type Order struct {
    ID              int        `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
    UserID          int        `gorm:"not null;column:user_id" json:"user_id"`
    Capacity        int        `gorm:"not null;column:capacity" json:"capacity"`
    LiftingCapacity int        `gorm:"not null;column:lifting_capacity" json:"lifting_capacity"`
    Type            *int       `gorm:"column:type" json:"type,omitempty"`                  
    ManagerID       *int       `gorm:"column:manager_id" json:"manager_id,omitempty"`       
    Status          int        `gorm:"not null;default:1;column:status" json:"status"`
    Price           *float64   `gorm:"type:decimal;column:price" json:"price,omitempty"`    
    PickupAddress   *string    `gorm:"column:pickup_address" json:"pickup_address,omitempty"`
    DeliveryAddress *string    `gorm:"column:delivery_address" json:"delivery_address,omitempty"`
    DateStart       *time.Time `gorm:"column:date_start" json:"date_start,omitempty"`
    DateEnd         *time.Time `gorm:"column:date_end" json:"date_end,omitempty"`
    PeriodFrom      *time.Time `gorm:"column:period_from" json:"period_from,omitempty"`
    PeriodTo        *time.Time `gorm:"column:period_to" json:"period_to,omitempty"`
    DateCreate      time.Time  `gorm:"not null;autoCreateTime;column:date_create" json:"date_create"`
    DateUpdate      *time.Time `gorm:"autoUpdateTime;column:date_update" json:"date_update,omitempty"` 
}