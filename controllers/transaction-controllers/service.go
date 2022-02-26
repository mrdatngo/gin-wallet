package transaction_controllers

import (
	model "github.com/mrdatngo/gin-wallet/models"
)

type Service interface {
	DepositListService(input *InputDepositList) (*[]model.EntityTransaction, string)
	DepositService(input *InputDeposit, userID string) (*model.EntityTransaction, string)
}

type service struct {
	repository Repository
}

func NewTransactionService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) DepositListService(input *InputDepositList) (*[]model.EntityTransaction, string) {

	resultTransactions, errResultTransactions := s.repository.DepositListRepository(input.UserID, input.WalletID)

	return resultTransactions, errResultTransactions
}

func (s *service) DepositService(input *InputDeposit, userID string) (*model.EntityTransaction, string) {

	transaction := model.EntityTransaction{
		Status:   0,
		Amount:   input.Amount,
		UserID:   userID,
		WalletID: input.WalletID,
		Type:     "Deposit",
	}

	resultCreateWallet, errCreateWallet := s.repository.DepositRepository(&transaction)

	return resultCreateWallet, errCreateWallet
}
