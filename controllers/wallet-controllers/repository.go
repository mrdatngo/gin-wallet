package wallet_controllers

import (
	_const "github.com/mrdatngo/gin-wallet/const"
	model "github.com/mrdatngo/gin-wallet/models"
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	ResultWalletsRepository(userId string) (*[]model.EntityWallet, string)
	CreateWalletRepository(input *model.EntityWallet) (*model.EntityWallet, string)
	DeleteWalletRepository(input *model.EntityWallet) (*model.EntityWallet, string)
}

type repository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ResultWalletsRepository(userId string) (*[]model.EntityWallet, string) {

	var wallets []model.EntityWallet
	db := r.db.Model(&wallets)
	errorCode := make(chan string, 1)

	resultsStudents := db.Debug().Select("*").Where("user_id = ? AND active = ?", userId, true).Find(&wallets)

	if resultsStudents.Error != nil {
		errorCode <- _const.RESULTS_WALLET_NOT_FOUND
		return &wallets, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &wallets, <-errorCode
}

func (r *repository) CreateWalletRepository(input *model.EntityWallet) (*model.EntityWallet, string) {

	var wallet model.EntityWallet

	wallet.UserID = input.UserID

	db := r.db.Model(&wallet)
	errorCode := make(chan string, 1)

	addNewWallet := db.Debug().Create(&wallet)
	db.Commit()

	if addNewWallet.Error != nil {
		log.Println(addNewWallet.Error)
		errorCode <- _const.CREATE_WALLET_FAILED
		return &wallet, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &wallet, ""
}

func (r *repository) DeleteWalletRepository(input *model.EntityWallet) (*model.EntityWallet, string) {

	var wallet model.EntityWallet
	db := r.db.Model(&wallet)
	errorCode := make(chan string, 1)

	checkWalletId := db.Debug().Select("*").Where("id = ? AND active = ?", input.ID, true).Find(&wallet)

	if checkWalletId.RowsAffected < 1 {
		errorCode <- _const.DELETE_WALLET_NOT_FOUND
		return &wallet, <-errorCode
	}

	if wallet.Balance > 0 {
		errorCode <- _const.DELETE_WALLET_BALANCE_NOT_ZERO
		return &wallet, <-errorCode
	}

	deleteWalletId := db.Debug().Select("*").Where("id = ?", input.ID).Find(&wallet).Update("active", false)

	if deleteWalletId.Error != nil {
		errorCode <- _const.DELETE_WALLET_FAILED
		return &wallet, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &wallet, <-errorCode
}
