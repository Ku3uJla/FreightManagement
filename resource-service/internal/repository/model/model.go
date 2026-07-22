package model

import "time"

type Driver struct {
	ID         int       `gorm:"type:int;primaryKey" json:"id"`
	UserID     int       `gorm:"unique" json:"user_id"`
	Status     int       `gorm:"type:integer" json:"status"`
	DateCreate time.Time `gorm:"autoCreateTime" json:"date_create"`
	DateUpdate time.Time `gorm:"autoUpdateTime" json:"date_update"`
}

type DriverCategory struct {
	ID         int       `gorm:"type:int;primaryKey" json:"id"`
	DriverID   int       `gorm:"type:int;not_null"`
	Category   string    `gorm:"size:10" json:"category"`
	DateCreate time.Time `gorm:"autoCreateTime" json:"date_create"`
	DateUpdate time.Time `gorm:"autoUpdateTime" json:"date_update"`
}

type Auto struct {
	ID               int       `gorm:"type:int;primaryKey" json:"id"`
	Status           int       `gorm:"type:integer" json:"status"`
	Capacity         int       `gorm:"type:integer;not_null"`
	LiftingCapacity  int       `gorm:"type:integer;not_null"`
	Number           string    `gorm:"size:10"`
	RequiredCategory string    `gorm:"size:10"`
	DateCreate       time.Time `gorm:"autoCreateTime" json:"date_create"`
	DateUpdate       time.Time `gorm:"autoUpdateTime" json:"date_update"`
}

type DriverAuto struct {
	ID         int `gorm:"type:int;primaryKey" json:"id"`
	DriverID   int `gorm:"unique" json:"driver_id"`
	OrderID    int `gorm:"type:int"`
	AutoID     int `gorm:"type:int"`
	Status     int `gorm:"type:integer" json:"status"`
	DateFinish time.Time
	DateCreate time.Time `gorm:"autoCreateTime" json:"date_create"`
	DateUpdate time.Time `gorm:"autoUpdateTime" json:"date_update"`
}
