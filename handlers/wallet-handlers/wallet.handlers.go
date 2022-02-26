package wallet_handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_const "github.com/mrdatngo/gin-wallet/const"
	auth_controllers "github.com/mrdatngo/gin-wallet/controllers/auth-controllers"
	wallet_controllers "github.com/mrdatngo/gin-wallet/controllers/wallet-controllers"
	util "github.com/mrdatngo/gin-wallet/utils"
	gpc "github.com/restuwahyu13/go-playground-converter"
	"net/http"
)

type handler struct {
	service     wallet_controllers.Service
	authService auth_controllers.Service
}

func NewWalletHandler(service wallet_controllers.Service, authService auth_controllers.Service) *handler {
	return &handler{service: service, authService: authService}
}

func (h *handler) ResultsWalletHandler(ctx *gin.Context) {

	var input wallet_controllers.InputResultsWallet
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "UserID",
				Message: "user_id is required on body",
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

	resultWallets, errResultWallets := h.service.ResultsWalletService(input.UserID)

	switch errResultWallets {

	case _const.RESULTS_WALLET_NOT_FOUND:
		util.APIResponse(ctx, "Wallets data is not exists", http.StatusConflict, http.MethodPost, nil)

	default:
		util.APIResponse(ctx, "Results wallets data successfully", http.StatusOK, http.MethodPost, resultWallets)
	}
}

func (h *handler) CreateWalletHandler(ctx *gin.Context) {

	var input wallet_controllers.InputCreateWallet
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "UserID",
				Message: "user_id is required on body",
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
	if user.RoleID != 1 {
		util.APIResponse(ctx, "Permission denied!", http.StatusForbidden, http.MethodPost, nil)
		return
	}

	_, errCreateWallet := h.service.CreateWalletService(&input)

	switch errCreateWallet {

	case _const.CREATE_WALLET_FAILED:
		util.APIResponse(ctx, "Create new wallet failed", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		util.APIResponse(ctx, "Create new wallet account successfully", http.StatusCreated, http.MethodPost, nil)
	}
}

func (h *handler) DeleteWalletHandler(ctx *gin.Context) {

	var input wallet_controllers.InputDeleteWallet
	input.ID = ctx.Param("id")

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "ID",
				Message: "id is required on param",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodDelete, errResponse)
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

	_, errDeleteWallet := h.service.DeleteWalletService(&input)

	switch errDeleteWallet {

	case _const.DELETE_WALLET_NOT_FOUND:
		util.APIResponse(ctx, "Wallet data is not exist or deleted", http.StatusForbidden, http.MethodDelete, nil)
		return

	case _const.DELETE_WALLET_FAILED:
		util.APIResponse(ctx, "Delete wallet data failed", http.StatusForbidden, http.MethodDelete, nil)
		return
	case _const.DELETE_WALLET_BALANCE_NOT_ZERO:
		util.APIResponse(ctx, "Can not delete wallet balance diff 0", http.StatusForbidden, http.MethodDelete, nil)
		return
	default:
		util.APIResponse(ctx, "Delete student data successfully", http.StatusOK, http.MethodDelete, nil)
	}
}
