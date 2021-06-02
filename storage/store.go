package storage

import "github.com/DAlconoiD/balance-service/models"

type Store interface {
	GetBalance(id int) (*models.Account, *models.CustomErr)
	GetTransactionHistory(accId int, sorting string, order string, page int) ([]models.Transaction, *models.CustomErr)
	UpdateBalance(request *models.ChangeBalanceRequest) (*models.Transaction, *models.CustomErr)
	MakeTransfer(request *models.TransferRequest) (*models.Transaction, *models.CustomErr)
}