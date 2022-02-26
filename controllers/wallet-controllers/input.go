package wallet_controllers

type InputResultsWallet struct {
	UserID string `json:"user_id" validate:"required"`
}

type InputCreateWallet struct {
	UserID string `json:"user_id" validate:"required"`
}

type InputDeleteWallet struct {
	ID string `validate:"required,uuid"`
}
