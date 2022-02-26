package transaction_controllers

import (
	_const "github.com/mrdatngo/gin-wallet/const"
	model "github.com/mrdatngo/gin-wallet/models"
	"gorm.io/gorm"
)

type Repository interface {
	DepositListRepository(userID string, walletID string) (*[]model.EntityTransaction, string)
	DepositRepository(input *model.EntityTransaction) (*model.EntityTransaction, string)
}

type repository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) DepositListRepository(userID string, walletID string) (*[]model.EntityTransaction, string) {

	var transactions []model.EntityTransaction
	db := r.db.Model(&transactions)
	errorCode := make(chan string, 1)

	resultTransactions := db.Debug().Select("*").Where("wallet_id = ?", walletID).Find(&transactions)

	if resultTransactions.Error != nil {
		errorCode <- _const.RESULTS_TRANSACTIONS_NOT_FOUND
		return nil, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &transactions, <-errorCode
}

func (r *repository) DepositRepository(input *model.EntityTransaction) (*model.EntityTransaction, string) {

	var wallet model.EntityWallet
	db := r.db.Model(&wallet)
	errorCode := make(chan string, 1)

	checkWalletId := db.Debug().Select("*").Where("id = ? AND active = ?", input.WalletID, true).Find(&wallet)

	if checkWalletId.RowsAffected < 1 {
		errorCode <- _const.WALLET_NOT_FOUND
		return nil, <-errorCode
	}

	var transaction model.EntityTransaction
	transaction = model.EntityTransaction{
		Description: "Deposit",
		Status:      0,
		Amount:      input.Amount,
		UserID:      input.UserID,
		WalletID:    input.WalletID,
	}

	tx := r.db.Begin()

	defer func() {
		if rcv := recover(); rcv != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return &transaction, err.Error()
	}

	if err := tx.Model(&transaction).Create(&transaction).Error; err != nil {
		tx.Rollback()
		return &transaction, err.Error()
	}

	if err := tx.Model(&wallet).Update("balance", wallet.Balance+transaction.Amount).Error; err != nil {
		tx.Rollback()
		return &transaction, err.Error()
	}

	if tx.Commit().Error != nil {
		return &transaction, tx.Commit().Error.Error()
	}

	return &transaction, "nil"
}
