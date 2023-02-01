package database

import (
	"finaltask/models"
	"finaltask/pkg/mysql"
	"fmt"
)

func Migrate() {
	err := mysql.DB.AutoMigrate(&models.User{}, &models.Films{}, &models.Genre{}, &models.Transaction{})

	if err != nil {
		fmt.Println(err)
		panic("Migration failed")
	}

	fmt.Println("Migration success")
}
