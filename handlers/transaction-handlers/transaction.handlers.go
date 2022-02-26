package transaction_handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_const "github.com/mrdatngo/gin-wallet/const"
	auth_controllers "github.com/mrdatngo/gin-wallet/controllers/auth-controllers"
	transaction_controllers "github.com/mrdatngo/gin-wallet/controllers/transaction-controllers"
	util "github.com/mrdatngo/gin-wallet/utils"
	gpc "github.com/restuwahyu13/go-playground-converter"
	"net/http"
)

type handler struct {
	service     transaction_controllers.Service
	authService auth_controllers.Service
}

func NewTransactionHandler(service transaction_controllers.Service, authService auth_controllers.Service) *handler {
	return &handler{service: service, authService: authService}
}

func (h *handler) DepositListHandler(ctx *gin.Context) {

	var input transaction_controllers.InputDepositList
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "UserID",
				Message: "user_id is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "WalletID",
				Message: "wallet_id is required on body",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	data, _ := ctx.Get("user")
	email := data.(jwt.MapClaims)["email"]
	if email == "" {
		util.APIResponse(ctx, "Something went wrong!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	user, errUser := h.authService.UserService(email.(string))
	switch errUser {
	case _const.USER_NOT_FOUND:
		util.APIResponse(ctx, "User login not found!", http.StatusForbidden, http.MethodPost, nil)
		return
	}
	if errUser != "" {
		util.APIResponse(ctx, errUser, http.StatusForbidden, http.MethodPost, nil)
		return
	}

	if !user.Active {
		util.APIResponse(ctx, "Your account is deactivate!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	userId := data.(jwt.MapClaims)["id"]

	if user.RoleID != 1 && userId != input.UserID {
		util.APIResponse(ctx, "Permission denied!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	resultDeposits, errResultDeposits := h.service.DepositListService(&input)

	switch errResultDeposits {

	case _const.RESULTS_TRANSACTIONS_NOT_FOUND:
		util.APIResponse(ctx, "Wallets transaction is not exists", http.StatusConflict, http.MethodPost, nil)

	default:
		util.APIResponse(ctx, "Results wallet transactions data successfully", http.StatusOK, http.MethodPost, resultDeposits)
	}
}

func (h *handler) DepositHandler(ctx *gin.Context) {

	var input transaction_controllers.InputDeposit
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "WalletID",
				Message: "wallet_id is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Amount",
				Message: "amount is required on body",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	if input.Amount <= 0 {
		util.APIResponse(ctx, "Amount must greater than 0!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	data, _ := ctx.Get("user")
	email := data.(jwt.MapClaims)["email"]
	if email == "" {
		util.APIResponse(ctx, "Something went wrong!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	user, errUser := h.authService.UserService(email.(string))
	switch errUser {
	case _const.USER_NOT_FOUND:
		util.APIResponse(ctx, "User login not found!", http.StatusForbidden, http.MethodPost, nil)
		return
	}
	if errUser != "" {
		util.APIResponse(ctx, errUser, http.StatusForbidden, http.MethodPost, nil)
		return
	}

	if !user.Active {
		util.APIResponse(ctx, "Your account is deactivate!", http.StatusForbidden, http.MethodPost, nil)
		return
	}
	if user.RoleID != 1 {
		util.APIResponse(ctx, "Permission denied!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	_, errDepositWallet := h.service.DepositService(&input, user.ID)

	switch errDepositWallet {

	case _const.WALLET_NOT_FOUND:
		util.APIResponse(ctx, "Wallet not found!", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		if errDepositWallet == "nil" || errDepositWallet == "" {
			util.APIResponse(ctx, "Transaction successfully", http.StatusCreated, http.MethodPost, nil)
		} else {
			util.APIResponse(ctx, errDepositWallet, http.StatusCreated, http.MethodPost, nil)
		}
	}
}
