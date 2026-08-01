package model

import "time"

type User struct {
	ID        int       `gorm:"type:int;primaryKey" json:"id"`
	Role      string    `gorm:"type:int" json:"role"`
	Phone     string    `gorm:"size:255;not null" json:"phone"`
	Email     string    `gorm:"size:255;not null" json:"email"`
	FullName  *string   `gorm:"size:255" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"dateCreate"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"dateUpdate"`
}
