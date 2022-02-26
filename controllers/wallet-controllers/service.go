package wallet_controllers

import (
	"github.com/google/uuid"
	model "github.com/mrdatngo/gin-wallet/models"
)

type Service interface {
	ResultsWalletService(userID string) (*[]model.EntityWallet, string)
	CreateWalletService(input *InputCreateWallet) (*model.EntityWallet, string)
	DeleteWalletService(input *InputDeleteWallet) (*model.EntityWallet, string)
}

type service struct {
	repository Repository
}

func NewWalletService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ResultsWalletService(userID string) (*[]model.EntityWallet, string) {

	resultWallets, errResultWallets := s.repository.ResultWalletsRepository(userID)

	return resultWallets, errResultWallets
}

func (s *service) CreateWalletService(input *InputCreateWallet) (*model.EntityWallet, string) {

	wallet := model.EntityWallet{
		ID:      uuid.New().String(),
		UserID:  input.UserID,
		Balance: 0,
	}

	resultCreateWallet, errCreateWallet := s.repository.CreateWalletRepository(&wallet)

	return resultCreateWallet, errCreateWallet
}

func (s *service) DeleteWalletService(input *InputDeleteWallet) (*model.EntityWallet, string) {

	wallet := model.EntityWallet{
		ID: input.ID,
	}

	resultDeleteWallet, errDeleteWallet := s.repository.DeleteWalletRepository(&wallet)

	return resultDeleteWallet, errDeleteWallet
}
