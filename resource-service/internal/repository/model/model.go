package model

import "time"

type Driver struct {
	ID        int       `gorm:"type:int;primaryKey" json:"id"`
	UserID    int       `gorm:"unique" json:"user_id"`
	Status    int       `gorm:"type:integer" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"date_create"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"date_update"`
}

type DriverCategory struct {
	ID        int       `gorm:"type:int;primaryKey" json:"id"`
	DriverID  int       `gorm:"type:int;not_null"`
	Category  string    `gorm:"size:10" json:"category"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"date_create"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"date_update"`
}

type Auto struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Status           int       `gorm:"type:integer" json:"status"`
	Capacity         int       `gorm:"type:integer;not null" json:"capacity" binding:"required"`
	LiftingCapacity  int       `gorm:"type:integer;not null" json:"lifting_capacity" binding:"required"`
	Number           string    `gorm:"size:10" json:"number" binding:"required"`
	RequiredCategory string    `gorm:"size:10" json:"required_category" binding:"required"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"date_create"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"date_update"`
}

type DriverAuto struct {
	ID         int `gorm:"type:int;primaryKey" json:"id"`
	DriverID   int `gorm:"unique" json:"driver_id"`
	OrderID    int `gorm:"type:int"`
	AutoID     int `gorm:"type:int"`
	Status     int `gorm:"type:integer" json:"status"`
	DateFinish time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"date_create"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"date_update"`
}
