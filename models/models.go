package models

import "time"

const (
	ErrorDefaultCode           = 0
	ErrorInsufficientFundsCode = 1

	InsufficientFundsMessage = "non_negative_balance"

	SortByTimeString = "by-time"
	SortBySumString  = "by-sum"

	OrderAscendingString  = "asc"
	OrderDescendingString = "desc"
)

type Account struct {
	ID      int `gorm:"primaryKey; column:account_id"`
	Balance float64
}

type Transaction struct {
	ID        int `gorm:"primaryKey; column:transaction_id"`
	AccountID int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Delta     float64
	Remaining float64
	Message   string
}

type CustomErr struct {
	Err       error
	ErrorCode int
}

type ChangeBalanceRequest struct {
	ID    int     `validate:"required,gt=0"`
	Delta float64 `validate:"required"`
}

type TransferRequest struct {
	ID1   int     `validate:"required,gt=0"`
	ID2   int     `validate:"required,nefield=ID1,gt=0"`
	Delta float64 `validate:"required,gt=0"`
}
