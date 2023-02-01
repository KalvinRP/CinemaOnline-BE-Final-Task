package repositories

import (
	"finaltask/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(user models.User) (models.User, error)
	Login(email string) (models.User, error)
	CheckAuth(ID int) (models.User, error)
	IsExist(user models.User) bool
	ResetPass(email, password string) (models.User, error)
}

func RepositoryAuth(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Register(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *repository) Login(email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "email=?", email).Error

	return user, err
}

func (r *repository) CheckAuth(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repository) IsExist(user models.User) bool {
	err := r.db.Where("email = ?", user.Email).First(&user).Error

	return err == nil
}

func (r *repository) ResetPass(email, password string) (models.User, error) {
	var user models.User
	err := r.db.Model(&user).Where("Email = ?", email).Update("Password", password).Error

	return user, err
}
