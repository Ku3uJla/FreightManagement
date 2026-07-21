package model

type User struct {
	User_ID  int    `gorm:"primaryKey;unique" json:"id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
	Password string `gorm:"size:100;not null" json:"password"`
}
