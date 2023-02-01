package models

import "time"

type Genre struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type GenreResponse struct {
	Name string
}

func (GenreResponse) TableName() string {
	return "genres"
}
