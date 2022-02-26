package auth_handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_const "github.com/mrdatngo/gin-wallet/const"
	auth_controllers "github.com/mrdatngo/gin-wallet/controllers/auth-controllers"
	util "github.com/mrdatngo/gin-wallet/utils"
	gpc "github.com/restuwahyu13/go-playground-converter"
	"github.com/sirupsen/logrus"
	"net/http"
)

type handler struct {
	service auth_controllers.Service
}

func NewAuthHandler(service auth_controllers.Service) *handler {
	return &handler{service: service}
}

func (h *handler) ActivationHandler(ctx *gin.Context) {
	var input auth_controllers.InputActivation

	token := ctx.Param("token")
	resultToken, errToken := util.VerifyToken(token, "JWT_SECRET")

	fmt.Println(token)
	fmt.Println(errToken)

	if errToken != nil {
		defer logrus.Error(errToken.Error())
		util.APIResponse(ctx, "Verified activation token failed", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	result := util.DecodeToken(resultToken)
	input.Email = result.Claims.Email
	input.Active = true

	_, errActivation := h.service.ActivationService(&input)

	switch errActivation {

	case "ACTIVATION_NOT_FOUND_404":
		util.APIResponse(ctx, "User account is not exist", http.StatusNotFound, http.MethodPost, nil)
		return

	case "ACTIVATION_ACTIVE_400":
		util.APIResponse(ctx, "User account hash been active please login", http.StatusBadRequest, http.MethodPost, nil)
		return

	case "ACTIVATION_ACCOUNT_FAILED_403":
		util.APIResponse(ctx, "Activation account failed", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		util.APIResponse(ctx, "Activation account success", http.StatusOK, http.MethodPost, nil)
	}
}

func (h *handler) ForgotHandler(ctx *gin.Context) {

	var input auth_controllers.InputForgot
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
		},
	}

	errResponse, errCount := util.GoValidator(input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	forgotResult, errForgot := h.service.ForgotService(&input)

	switch errForgot {

	case "FORGOT_NOT_FOUD_404":
		util.APIResponse(ctx, "Email is not never registered", http.StatusNotFound, http.MethodPost, nil)
		return

	case "FORGOT_NOT_ACTIVE_403":
		util.APIResponse(ctx, "User account is not active", http.StatusForbidden, http.MethodPost, nil)
		return

	case "FORGOT_PASSWORD_FAILED_403":
		util.APIResponse(ctx, "Forgot password failed", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		accessTokenData := map[string]interface{}{"id": forgotResult.ID, "email": forgotResult.Email}
		accessToken, errToken := util.Sign(accessTokenData, "JWT_SECRET", 5)

		if errToken != nil {
			defer logrus.Error(errToken.Error())
			util.APIResponse(ctx, "Generate accessToken failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		_, errorEmail := util.SendGridMail(forgotResult.Fullname, forgotResult.Email, "Reset Password", "template_reset", accessToken)

		if errorEmail != nil {
			util.APIResponse(ctx, "Sending email reset password failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		util.APIResponse(ctx, "Forgot password successfully", http.StatusOK, http.MethodPost, nil)
	}
}

func (h *handler) LoginHandler(ctx *gin.Context) {

	var input auth_controllers.InputLogin
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Password",
				Message: "password is required on body",
			},
		},
	}

	errResponse, errCount := util.GoValidator(&input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	resultLogin, errLogin := h.service.LoginService(&input)

	switch errLogin {

	case "LOGIN_NOT_FOUND_404":
		util.APIResponse(ctx, "User account is not registered", http.StatusNotFound, http.MethodPost, nil)
		return

	case "LOGIN_NOT_ACTIVE_403":
		util.APIResponse(ctx, "User account is not active", http.StatusForbidden, http.MethodPost, nil)
		return

	case "LOGIN_WRONG_PASSWORD_403":
		util.APIResponse(ctx, "Username or password is wrong", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		accessTokenData := map[string]interface{}{"id": resultLogin.ID, "email": resultLogin.Email}
		accessToken, errToken := util.Sign(accessTokenData, "JWT_SECRET", 24*60*1)

		if errToken != nil {
			defer logrus.Error(errToken.Error())
			util.APIResponse(ctx, "Generate accessToken failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		util.APIResponse(ctx, "Login successfully", http.StatusOK, http.MethodPost, map[string]string{"accessToken": accessToken})
	}
}

func (h *handler) RegisterHandler(ctx *gin.Context) {

	var input auth_controllers.InputRegister
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Fullname",
				Message: "fullname is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "lowercase",
				Field:   "Fullname",
				Message: "fullname must be using lowercase",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Password",
				Message: "password is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "gte",
				Field:   "Password",
				Message: "password minimum must be 8 character",
			},
		},
	}

	errResponse, errCount := util.GoValidator(input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	resultRegister, errRegister := h.service.RegisterService(&input)

	switch errRegister {

	case "REGISTER_CONFLICT_409":
		util.APIResponse(ctx, "Email already exist", http.StatusConflict, http.MethodPost, nil)
		return

	case "REGISTER_FAILED_403":
		util.APIResponse(ctx, "Register new account failed", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		accessTokenData := map[string]interface{}{"id": resultRegister.ID, "email": resultRegister.Email}
		accessToken, errToken := util.Sign(accessTokenData, "JWT_SECRET", 60)

		println("###" + accessToken + "###")

		if errToken != nil {
			defer logrus.Error(errToken.Error())
			util.APIResponse(ctx, "Generate accessToken failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		_, errSendMail := util.SendGridMail(resultRegister.Fullname, resultRegister.Email, "Activation Account", "template_register", accessToken)

		if errSendMail != nil {
			defer logrus.Error(errSendMail.Error())
			util.APIResponse(ctx, "Sending email activation failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		data := map[string]interface{}{}
		if util.GodotEnv("GO_ENV") != "production" {
			data["token"] = accessToken
		}
		util.APIResponse(ctx, "Register new account successfully", http.StatusCreated, http.MethodPost, data)
	}
}

func (h *handler) ResendHandler(ctx *gin.Context) {

	var input auth_controllers.InputResend
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
		},
	}

	errResponse, errCount := util.GoValidator(input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	resendResult, errResend := h.service.ResendService(&input)

	switch errResend {

	case "RESEND_NOT_FOUD_404":
		util.APIResponse(ctx, "Email is not never registered", http.StatusNotFound, http.MethodPost, nil)
		return

	case "RESEND_ACTIVE_403":
		util.APIResponse(ctx, "User account hash been active", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		accessTokenData := map[string]interface{}{"id": resendResult.ID, "email": resendResult.Email}
		accessToken, errToken := util.Sign(accessTokenData, "JWT_SECRET", 5)

		if errToken != nil {
			defer logrus.Error(errToken.Error())
			util.APIResponse(ctx, "Generate accessToken failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		_, errorSendEmail := util.SendGridMail(resendResult.Fullname, resendResult.Email, "Resend New Activation", "template_resend", accessToken)

		if errorSendEmail != nil {
			defer logrus.Error(errorSendEmail.Error())
			util.APIResponse(ctx, "Sending email resend activation failed", http.StatusBadRequest, http.MethodPost, nil)
			return
		}

		util.APIResponse(ctx, "Resend new activation token successfully", http.StatusOK, http.MethodPost, nil)
	}
}

func (h *handler) ResetHandler(ctx *gin.Context) {

	var input auth_controllers.InputReset
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Password",
				Message: "password is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "gte",
				Field:   "Password",
				Message: "password minimum must be 8 character",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Cpassword",
				Message: "cpassword is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "gte",
				Field:   "Cpassword",
				Message: "cpassword minimum must be 8 character",
			},
		},
	}

	errResponse, errCount := util.GoValidator(input, config.Options)

	if errCount > 0 {
		util.ValidatorErrorResponse(ctx, http.StatusBadRequest, http.MethodPost, errResponse)
		return
	}

	token := ctx.Param("token")
	resultToken, errToken := util.VerifyToken(token, "JWT_SECRET")

	if errToken != nil {
		defer logrus.Error(errToken.Error())
		util.APIResponse(ctx, "Verified activation token failed", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	if input.Cpassword != input.Password {
		util.APIResponse(ctx, "Confirm Password is not match with Password", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	result := util.DecodeToken(resultToken)
	input.Email = result.Claims.Email
	input.Active = true

	_, errReset := h.service.ResetService(&input)

	switch errReset {

	case "RESET_NOT_FOUND_404":
		util.APIResponse(ctx, "User account is not exist", http.StatusNotFound, http.MethodPost, nil)
		return

	case "ACCOUNT_NOT_ACTIVE_403":
		util.APIResponse(ctx, "User account is not active", http.StatusForbidden, http.MethodPost, nil)
		return

	case "RESET_PASSWORD_FAILED_403":
		util.APIResponse(ctx, "Change new password failed", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		util.APIResponse(ctx, "Change new password successfully", http.StatusOK, http.MethodPost, nil)
	}
}

func (h *handler) UpdateEmailHandler(ctx *gin.Context) {

	var input auth_controllers.InputUpdateEmail
	ctx.ShouldBindJSON(&input)

	config := gpc.ErrorConfig{
		Options: []gpc.ErrorMetaConfig{
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "Email",
				Message: "email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "Email",
				Message: "email format is not valid",
			},
			gpc.ErrorMetaConfig{
				Tag:     "required",
				Field:   "NewEmail",
				Message: "new_email is required on body",
			},
			gpc.ErrorMetaConfig{
				Tag:     "email",
				Field:   "NewEmail",
				Message: "new_email format is not valid",
			},
		},
	}

	errResponse, errCount := util.GoValidator(input, config.Options)

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

	user, errUser := h.service.UserService(email.(string))
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

	_, errUpdate := h.service.UpdateEmailService(&input)

	switch errUpdate {

	case _const.USER_NOT_FOUND:
		util.APIResponse(ctx, "User account is not found!", http.StatusForbidden, http.MethodPost, nil)
		return

	case _const.ACCOUNT_DEACTIVE:
		util.APIResponse(ctx, "User account is not active!", http.StatusForbidden, http.MethodPost, nil)
		return

	case _const.UPDATE_EMAIL_FAILED:
		util.APIResponse(ctx, "Update email failed!", http.StatusForbidden, http.MethodPost, nil)
		return

	default:
		util.APIResponse(ctx, "Update email successfully", http.StatusOK, http.MethodPost, nil)
	}
}
