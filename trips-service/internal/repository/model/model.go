package model

import (
	"time"
)

type Trip struct {
	ID         int        `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	DriverID   *int       `gorm:"column:driver_id" json:"driver_id,omitempty"` // может быть NULL
	OrderID    int        `gorm:"not null;column:order_id" json:"order_id"`    // обязательно
	AutoID     *int       `gorm:"column:auto_id" json:"auto_id,omitempty"`     // может быть NULL
	Status     string     `gorm:"type:varchar(255);not null;default:'0';column:status" json:"status"`
	DateFinish *time.Time `gorm:"column:date_finish" json:"date_finish,omitempty"` // может быть NULL
	CreatedAt  time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
