package model

import "time"

type Auth struct {
	ID        int       `gorm:"type:int;primaryKey" json:"id"`
	Login     string    `gorm:"size:255;not null;unique" json:"login"`
	Email     string    `gorm:"size:255;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"dateCreate"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"dateUpdate"`
}
