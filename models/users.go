package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Name      string    `json:"name" form:"name" gorm:"type: varchar(255)"`
	Email     string    `json:"email" form:"email" gorm:"type: varchar(255)"`
	Password  string    `json:"password" form:"password" gorm:"type: varchar(255)"`
	Role      string    `gorm:"type:varchar(5)"`
	Image     string    `json:"image" form:"image" gorm:"type: varchar(255)"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type UserDetail struct {
	ID    int
	Name  string
	Email string
	Image string
}

func (UserDetail) TableName() string {
	return "users"
}

type AuthResponse struct {
	ID    int
	Name  string
	Email string
}

func (AuthResponse) TableName() string {
	return "users"
}
