package models

import "time"

type Transaction struct {
	ID        string     `json:"id" gorm:"primary_key:auto_increment"`
	FilmsID   int        `json:"films_id"`
	Films     Films      `json:"films" gorm:"foreignKey:FilmsID"`
	UsersID   int        `json:"users_id"`
	Users     UserDetail `json:"users" gorm:"foreignKey:UsersID"`
	Status    string     `json:"status" gorm:"type: varchar(20)"`
	TempToken string     `json:"temptoken" gorm:"type: varchar(255)"`
	CreatedAt time.Time  `json:"buydate"`
	UpdatedAt time.Time  `json:"-"`
}

type TransactionResponse struct {
	ID        int        `json:"id"`
	Status    string     `json:"status"`
	FilmsID   int        `json:"films_id"`
	Films     Films      `json:"films" gorm:"foreignKey: films_id"`
	UsersID   int        `json:"users_id"`
	Users     UserDetail `json:"users" gorm:"foreignKey: users_id"`
	TempToken string     `json:"temptoken" gorm:"type: varchar(255)"`
	CreatedAt time.Time  `json:"buydate"`
}

func (TransactionResponse) TableName() string {
	return "profiles"
}
