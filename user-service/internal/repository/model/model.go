package model

import "time"

type User struct {
	ID         int       `gorm:"type:int;primaryKey" json:"id"`
	Role       int       `gorm:"type:int" json:"role"`
	Phone      string    `gorm:"size:255;not null" json:"phone"`
	Email      string    `gorm:"size:255;not null" json:"email"`
	Password   string    `gorm:"size:255;not null" json:"-"`
	Name       string    `gorm:"size:255" json:"name"`
	Login      string    `gorm:"size:255;not null;unique" json:"login"`
	DateCreate time.Time `gorm:"autoCreateTime" json:"dateCreate"`
	DateUpdate time.Time `gorm:"autoUpdateTime" json:"dateUpdate"`
}
