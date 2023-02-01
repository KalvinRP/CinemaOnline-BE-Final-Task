package filmsdto

type FilmsRequest struct {
	Title   string `json:"title" form:"title" gorm:"type: varchar(255)" validate:"required"`
	Image   string `json:"image" form:"image" gorm:"type: varchar(255)"`
	GenreID int    `json:"genre_id" form:"genre_id" gorm:"type: int" validate:"required"`
	Price   int    `json:"price" form:"price" gorm:"type: int" validate:"required"`
	Desc    string `json:"desc" form:"desc" gorm:"type: text" validate:"required"`
	Status  string `json:"status" form:"isbo" gorm:"type: varchar(10)" validate:"required"`
	YTID    string `json:"ytid" form:"ytid" gorm:"type: varchar(10)" validate:"required"`
	FullUrl string `json:"full_url" form:"full_url" gorm:"type: varchar(255)" validate:"required"`
}
