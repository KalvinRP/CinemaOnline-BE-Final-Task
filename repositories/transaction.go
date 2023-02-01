package repositories

import (
	"finaltask/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	GetOneTrans(ID string) (models.Transaction, error)
	GetAllTrans() ([]models.Transaction, error)
	GetTransbyStatus(Status string) (models.Transaction, error)
	AddTrans(transaction models.Transaction) (models.Transaction, error)
	UpdateTrans(status string, ID models.Transaction) error
	UserHistory(ID int) ([]models.Transaction, error)
	UserFilms(ID int) ([]models.Transaction, error)
	InputToken(token, transactionId string) error
	IsBought(ID int) bool
}

func RepoTrans(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetOneTrans(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload(clause.Associations).Preload("Films.Genre").First(&transaction, "id = ?", ID).Error

	return transaction, err
}

func (r *repository) GetAllTrans() ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Order("created_at DESC").Preload(clause.Associations).Preload("Films.Genre").Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransbyStatus(Status string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload(clause.Associations).Preload("Films.Genre").Find(&transaction, "status = ?", Status).Error

	return transaction, err
}

func (r *repository) AddTrans(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error

	return transaction, err
}

func (r *repository) UpdateTrans(status string, transaction models.Transaction) error {
	if status != transaction.Status && status == "success" {
		var film models.Films
		transaction.TempToken = ""

		r.db.Model(&transaction).Where("ID = ?", transaction.ID).Update("TempToken", nil)
		r.db.Model(&film).Where("ID = ?", transaction.Films.ID).Update("Sold", transaction.Films.Sold+1)
	}

	if status != transaction.Status && status == "failed" {
		transaction.TempToken = ""

		r.db.Model(&transaction).Where("ID = ?", transaction.ID).Update("TempToken", nil)
	}

	transaction.Status = status

	err := r.db.Save(&transaction).Error

	return err
}

func (r *repository) InputToken(token string, transactionId string) error {
	var trans models.Transaction
	err := r.db.Model(&trans).Where("ID = ?", transactionId).Update("TempToken", token).Error

	return err
}

func (r *repository) UserHistory(ID int) ([]models.Transaction, error) {
	var history []models.Transaction
	err := r.db.Preload(clause.Associations).Preload("Films.Genre").Where("users_id = ?", ID).Find(&history).Error

	return history, err
}

func (r *repository) UserFilms(ID int) ([]models.Transaction, error) {
	var history []models.Transaction
	err := r.db.Preload(clause.Associations).Preload("Films.Genre").Where("users_id = ? AND status = ?", ID, "success").Find(&history).Error

	return history, err
}

func (r *repository) IsBought(ID int) bool {
	var trans models.Transaction
	err := r.db.Where("films_id = ?", ID).First(&trans).Error

	return err == nil
}
