package repositories

import (
	"finaltask/models"

	"gorm.io/gorm"
)

type GenreRepository interface {
	AddGenre(genre models.Genre) (models.Genre, error)
	FindGenre(ID int) (models.Genre, error)
	GetGenre() ([]models.Genre, error)
	EditGenre(genre models.Genre, ID int) (models.Genre, error)
	DeleteGenre(genre models.Genre, ID int) (models.Genre, error)
}

func RepoGenre(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddGenre(genre models.Genre) (models.Genre, error) {
	err := r.db.Create(&genre).Error

	return genre, err
}

func (r *repository) FindGenre(ID int) (models.Genre, error) {
	var previousGenre models.Genre
	err := r.db.First(&previousGenre, ID).Error

	return previousGenre, err
}

func (r *repository) GetGenre() ([]models.Genre, error) {
	var genre []models.Genre
	err := r.db.Order("name").Find(&genre).Error

	return genre, err
}

func (r *repository) EditGenre(genre models.Genre, ID int) (models.Genre, error) {
	err := r.db.Save(&genre).Error

	return genre, err
}

func (r *repository) DeleteGenre(genre models.Genre, ID int) (models.Genre, error) {
	err := r.db.Delete(&genre).Error

	return genre, err
}
