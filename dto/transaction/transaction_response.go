package transactiondto

import (
	"finaltask/models"
)

type TransactionResponse struct {
	Status string
	Film   models.Films
	User   models.UserDetail
}
