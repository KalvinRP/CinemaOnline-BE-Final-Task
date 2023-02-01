package transactiondto

type TransactionRequest struct {
	FilmsID int    `json:"films_id"`
	UsersID int    `json:"users_id"`
	Status  string `json:"status" gorm:"type: varchar(20)"`
}

type UpdateTransactionRequest struct {
	Status string
}
