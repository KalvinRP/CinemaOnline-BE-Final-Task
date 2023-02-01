package models

import "time"

type Films struct {
	ID        int       `json:"id" gorm:"primary_key:auto_increment"`
	Title     string    `json:"title" form:"title" gorm:"type: varchar(255)"`
	Desc      string    `json:"desc" form:"desc" gorm:"type: text"`
	Price     int       `json:"price" form:"price" gorm:"type: int"`
	Sold      int       `json:"sold" form:"sold" gorm:"type: int"`
	Image     string    `json:"image" form:"image" gorm:"type: varchar(255)"`
	YTID      string    `json:"ytid" form:"ytid" gorm:"type: varchar(11)"`
	FullUrl   string    `json:"full_url" form:"full_url" gorm:"type: varchar(255)"`
	Token     string    `json:"token" form:"token" gorm:"type: varchar(255)"`
	Status    string    `json:"status" form:"isbo" gorm:"type: varchar(10)"` //ad, new, (favorite), regular
	Genre     Genre     `json:"genre" gorm:"foreignKey: GenreID"`
	GenreID   int       `json:"genre_id" form:"genre_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type FilmsResponse struct {
	Title   string
	Genre   Genre
	GenreID int
	Price   int
	Image   string
	Status  string
	Sold    int
}
