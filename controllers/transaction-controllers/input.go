package transaction_controllers

type InputDeposit struct {
	WalletID string `json:"wallet_id" validate:"required"`
	Amount   int64  `json:"amount" validate:"required"`
}

type InputDepositList struct {
	WalletID string `json:"wallet_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
}
