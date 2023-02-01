package repositories

import (
	"finaltask/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAcc(ID int) (models.User, error)
	EditAcc(user models.User) (models.User, error)
	DeleteAcc(user models.User, ID int) (models.User, error)
}

func RepoUser(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetAcc(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repository) EditAcc(user models.User) (models.User, error) {
	err := r.db.Model(&user).Updates(user).Error

	return user, err
}

func (r *repository) DeleteAcc(user models.User, ID int) (models.User, error) {
	err := r.db.Delete(&user).Error

	return user, err
}
