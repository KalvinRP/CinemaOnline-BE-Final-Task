package repositories

import (
	"finaltask/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilmsRepository interface {
	AddFilms(films models.Films) (models.Films, error)
	GetFilms() ([]models.Films, error)
	GetTopFilms() ([]models.Films, error)
	FindOneFilm(ID int) (models.Films, error)
	FindAFilm(ID, User int) (models.Films, error)
	EditFilms(films models.Films, ID int) (models.Films, error)
	DeleteFilms(films models.Films, ID int) (models.Films, error)
}

func RepoFilms(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddFilms(films models.Films) (models.Films, error) {
	err := r.db.Create(&films).Error

	return films, err
}

func (r *repository) GetTopFilms() ([]models.Films, error) {
	var films []models.Films
	// err := r.db.Preload(clause.Associations).Find(&films, []string{"ad", "new"}).Error
	subQuery := r.db.Table("films").Select("MAX(sold)")
	err := r.db.Preload(clause.Associations).Where("status IN ?", []string{"ad", "new"}).Or("sold = (?)", subQuery).Find(&films).Error

	return films, err
}

func (r *repository) FindOneFilm(ID int) (models.Films, error) {
	var film models.Films
	err := r.db.Preload("Genre").First(&film, ID).Error

	return film, err
}

func (r *repository) FindAFilm(ID, User int) (models.Films, error) {
	var film models.Films
	err := r.db.Preload("Genre").First(&film, ID).Error

	var trans models.Transaction
	error := r.db.Where("films_id = ? AND users_id = ? ", ID, User).First(&trans).Error
	if error == nil && trans.Status == "success" {
		film.Price = 0
	}
	if trans.Status == "pending" {
		film.Token = trans.TempToken
	}
	return film, err
}

func (r *repository) GetFilms() ([]models.Films, error) {
	var films []models.Films
	err := r.db.Preload(clause.Associations).Find(&films).Error

	return films, err
}

func (r *repository) EditFilms(films models.Films, ID int) (models.Films, error) {
	err := r.db.Model(&films).Updates(films).Error

	return films, err
}

func (r *repository) DeleteFilms(films models.Films, ID int) (models.Films, error) {
	var film models.Films
	err := r.db.Delete(&films).Error

	return film, err
}
