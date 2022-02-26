package auth_controllers

import (
	_const "github.com/mrdatngo/gin-wallet/const"
	model "github.com/mrdatngo/gin-wallet/models"
	util "github.com/mrdatngo/gin-wallet/utils"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	ActivationRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	ForgotRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	LoginRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	RegisterRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	ResendRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	ResetRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	UserRepository(input *model.EntityUsers) (*model.EntityUsers, string)
	UpdateEmailRepository(oldEmail string, newEmail string, metaData string) (*model.EntityUsers, string)
}

type repository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ActivationRepository(input *model.EntityUsers) (*model.EntityUsers, string) {

	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = input.Email

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- "ACTIVATION_NOT_FOUND_404"
		return &users, <-errorCode
	}

	db.Debug().Select("Active").Where("activation = ?", input.Active).Take(&users)

	if users.Active {
		errorCode <- "ACTIVATION_ACTIVE_400"
		return &users, <-errorCode
	}

	users.Active = input.Active
	users.UpdatedAt = time.Now().Local()

	updateActivation := db.Debug().Select("active", "updated_at").Where("email = ?", input.Email).Updates(users)

	if updateActivation.Error != nil {
		errorCode <- "ACTIVATION_ACCOUNT_FAILED_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}

func (r *repository) ForgotRepository(input *model.EntityUsers) (*model.EntityUsers, string) {

	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = input.Email
	users.Password = util.HashPassword(util.RandStringBytes(20))

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- "FORGOT_NOT_FOUD_404"
		return &users, <-errorCode
	}

	if !users.Active {
		errorCode <- "FORGOT_NOT_ACTIVE_403"
		return &users, <-errorCode
	}

	changePassword := db.Debug().Select("password", "updated_at").Where("email = ?", input.Email).Updates(users)

	if changePassword.Error != nil {
		errorCode <- "FORGOT_PASSWORD_FAILED_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}

func (r *repository) LoginRepository(input *model.EntityUsers) (*model.EntityUsers, string) {

	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = input.Email
	users.Password = input.Password

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- "LOGIN_NOT_FOUND_404"
		return &users, <-errorCode
	}

	if !users.Active {
		errorCode <- "LOGIN_NOT_ACTIVE_403"
		return &users, <-errorCode
	}

	comparePassword := util.ComparePassword(users.Password, input.Password)

	if comparePassword != nil {
		errorCode <- "LOGIN_WRONG_PASSWORD_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}

func (r *repository) RegisterRepository(input *model.EntityUsers) (*model.EntityUsers, string) {

	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected > 0 {
		errorCode <- "REGISTER_CONFLICT_409"
		return &users, <-errorCode
	}

	users.Fullname = input.Fullname
	users.Email = input.Email
	users.Password = input.Password

	addNewUser := db.Debug().Create(&users)
	db.Commit()

	if addNewUser.Error != nil {
		errorCode <- "REGISTER_FAILED_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}

func (r *repository) ResendRepository(input *model.EntityUsers) (*model.EntityUsers, string) {

	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = input.Email

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- "RESEND_NOT_FOUD_404"
		return &users, <-errorCode
	}

	if users.Active {
		errorCode <- "RESEND_ACTIVE_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}

func (r *repository) ResetRepository(input *model.EntityUsers) (*model.EntityUsers, string) {
	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = input.Email
	users.Password = input.Password
	users.Active = input.Active

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- "RESET_NOT_FOUND_404"
		return &users, <-errorCode
	}

	if !users.Active {
		errorCode <- "ACCOUNT_NOT_ACTIVE_403"
		return &users, <-errorCode
	}

	users.Password = util.HashPassword(input.Password)
	users.UpdatedAt = time.Now().Local()

	updateNewPassword := db.Debug().Select("password", "update_at").Where("email = ?", input.Email).Updates(users)

	if updateNewPassword.Error != nil {
		errorCode <- "RESET_PASSWORD_FAILED_403"
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}

func (r *repository) UserRepository(input *model.EntityUsers) (*model.EntityUsers, string) {
	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	checkUserAccount := db.Debug().Select("*").Where("email = ?", input.Email).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- _const.USER_NOT_FOUND
		return &users, <-errorCode
	}

	if checkUserAccount.Error != nil {
		return nil, checkUserAccount.Error.Error()
	}
	return &users, ""
}

func (r *repository) UpdateEmailRepository(oldEmail string, newEmail string, metaData string) (*model.EntityUsers, string) {

	var users model.EntityUsers
	db := r.db.Model(&users)
	errorCode := make(chan string, 1)

	users.Email = newEmail

	checkUserAccount := db.Debug().Select("*").Where("email = ?", oldEmail).Find(&users)

	if checkUserAccount.RowsAffected < 1 {
		errorCode <- _const.USER_NOT_FOUND
		return &users, <-errorCode
	}

	db.Debug().Select("Active").Where("active = ?", true).Take(&users)

	if !users.Active {
		errorCode <- _const.ACCOUNT_DEACTIVE
		return &users, <-errorCode
	}

	users.Email = newEmail
	users.MetaData = metaData
	users.UpdatedAt = time.Now().Local()

	updateEmail := db.Debug().Select("*").Where("email = ?", oldEmail).Updates(users)

	if updateEmail.Error != nil {
		errorCode <- _const.UPDATE_EMAIL_FAILED
		return &users, <-errorCode
	} else {
		errorCode <- "nil"
	}

	return &users, <-errorCode
}
