package storage

import (
	"fmt"
	"github.com/DAlconoiD/balance-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

//Database represents a real database
type Database struct {
	Db         *gorm.DB
	ConnString string
	PaginationNum int
}

//Open establishes a connection to database
func (db *Database) Open() error {
	var err error
	db.Db, err = gorm.Open(postgres.Open(db.ConnString), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

//GetBalance returns account with id=id
func (db *Database) GetBalance(id int) (*models.Account, *models.CustomErr) {
	var account = &models.Account{}
	result := db.Db.First(account, id)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &models.Account{ID: id, Balance: 0}, nil
	} else if result.Error != nil {
		return nil, &models.CustomErr{Err: fmt.Errorf("GetBalance: %v", result.Error), ErrorCode: models.ErrorDefaultCode}
	}

	return account, nil
}

//GetTransactionHistory returns transaction history sorted by time/sum asc/desc; supports pagination
func (db *Database) GetTransactionHistory(accId int, sorting string, order string, page int) ([]models.Transaction, *models.CustomErr) {
	history := make([]models.Transaction, 0, 0)

	query := db.Db.Where("account_id = ?", accId)
	var sortStr string
	switch sorting {
	case models.SortByTimeString:
		sortStr = "created_at "
	case models.SortBySumString:
		sortStr = "delta "
	}
	sortStr += order
	query.Order(sortStr)

	if page > 0 {
		query.Limit(db.PaginationNum).Offset((page-1)*db.PaginationNum)
	}

	result := query.Find(&history)
	if result.Error != nil {
		return nil, &models.CustomErr{
			Err:       result.Error,
			ErrorCode: models.ErrorDefaultCode,
		}
	}

	return history, nil
}

//UpdateBalance changes account balance
func (db *Database) UpdateBalance(request *models.ChangeBalanceRequest) (*models.Transaction, *models.CustomErr) {
	tx := db.Db.Begin()
	now := time.Now()
	account, err := updOrCreateAccBalance(tx, request.ID, request.Delta)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction := &models.Transaction{
		AccountID: account.ID,
		CreatedAt: now,
		Delta:     request.Delta,
		Remaining: account.Balance,
		Message:   fmt.Sprintf("Account [%v]: balance changed by [%.2f], [%.2f] remaining", account.ID, request.Delta, account.Balance),
	}

	if err = writeTransaction(tx, transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return transaction, nil
}

//MakeTransfer makes transfer between accounts
func (db *Database) MakeTransfer(request *models.TransferRequest) (*models.Transaction, *models.CustomErr) {
	tx := db.Db.Begin()
	now := time.Now()

	account1, err := updOrCreateAccBalance(tx, request.ID1, -request.Delta)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	account2, err := updOrCreateAccBalance(tx, request.ID2, request.Delta)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	transaction1 := &models.Transaction{
		AccountID: account1.ID,
		CreatedAt: now,
		Delta:     -request.Delta,
		Remaining: account1.Balance,
		Message: fmt.Sprintf("Transfer from account [%v] to account [%v]: balance changed by [%.2f], [%.2f] remaining",
			account1.ID, account2.ID, -request.Delta, account1.Balance),
	}
	transaction2 := &models.Transaction{
		AccountID: account2.ID,
		CreatedAt: now,
		Delta:     request.Delta,
		Remaining: account2.Balance,
		Message: fmt.Sprintf("Transfer from account [%v] to account [%v]: balance changed by [%.2f], [%.2f] remaining",
			account1.ID, account2.ID, request.Delta, account2.Balance),
	}

	if err = writeTransaction(tx, transaction1); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = writeTransaction(tx, transaction2); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return transaction1, err
}

func updOrCreateAccBalance(tx *gorm.DB, id int, delta float64) (*models.Account, *models.CustomErr) {
	result := tx.Model(&models.Account{ID: id}).UpdateColumn("balance", gorm.Expr("balance + ?", delta))
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), models.InsufficientFundsMessage) {
			return nil, &models.CustomErr{
				Err:       fmt.Errorf("insuffisient funds on account [%v]", id),
				ErrorCode: models.ErrorInsufficientFundsCode,
			}
		}
		return nil, &models.CustomErr{Err: result.Error, ErrorCode: models.ErrorDefaultCode}
	}
	if result.RowsAffected == 0 {
		//create account if delta > 0
		if delta >= 0 {
			account := &models.Account{ID: id, Balance: delta}
			result = tx.Create(account)
		} else {
			return nil, &models.CustomErr{
				Err:       fmt.Errorf("insuffisient funds on account [%v]", id),
				ErrorCode: models.ErrorInsufficientFundsCode,
			}
		}
	}
	fmt.Printf("UPDATE BALANCE: rows affected = [%v]", result.RowsAffected)
	account := &models.Account{}
	tx.First(account, id)

	return account, nil
}

func writeTransaction(tx *gorm.DB, transaction *models.Transaction) *models.CustomErr {
	result := tx.Create(transaction)
	if result.Error != nil {
		return &models.CustomErr{
			Err:       result.Error,
			ErrorCode: models.ErrorDefaultCode,
		}
	}
	fmt.Printf("WRITE TRANSACTION: rows affected = [%v]", result.RowsAffected)
	return nil
}

