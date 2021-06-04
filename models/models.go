package models

import "time"

const (
	//custom error codes
	ErrorDefaultCode           = 0
	ErrorInsufficientFundsCode = 1

	//name of database constraint
	InsufficientFundsMessage = "non_negative_balance"

	//valid URL query "sorted" param values
	SortByTimeString = "by-time"
	SortBySumString  = "by-sum"

	//valid URL query "order" param values
	OrderAscendingString  = "asc"
	OrderDescendingString = "desc"
)

//Account - account model
type Account struct {
	ID      int `gorm:"primaryKey; column:account_id"`
	Balance float64
}

//Transaction - transaction model
type Transaction struct {
	ID        int `gorm:"primaryKey; column:transaction_id"`
	AccountID int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Delta     float64
	Remaining float64
	Message   string
}

//CustomErr - custom error model
type CustomErr struct {
	Err       error
	ErrorCode int
}

//ChangeBalanceRequest is a model which handleChangeBalance expects
type ChangeBalanceRequest struct {
	ID    int     `validate:"required,gt=0"`
	Delta float64 `validate:"required"`
}

//TransferRequest is a model which handleTransfer expects
type TransferRequest struct {
	ID1   int     `validate:"required,gt=0"`
	ID2   int     `validate:"required,nefield=ID1,gt=0"`
	Delta float64 `validate:"required,gt=0"`
}
