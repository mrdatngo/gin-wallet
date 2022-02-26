package auth_controllers

import (
	"fmt"
	model "github.com/mrdatngo/gin-wallet/models"
	"time"
)

type Service interface {
	ActivationService(input *InputActivation) (*model.EntityUsers, string)
	ForgotService(input *InputForgot) (*model.EntityUsers, string)
	LoginService(input *InputLogin) (*model.EntityUsers, string)
	RegisterService(input *InputRegister) (*model.EntityUsers, string)
	ResendService(input *InputResend) (*model.EntityUsers, string)
	ResetService(input *InputReset) (*model.EntityUsers, string)
	UserService(email string) (*model.EntityUsers, string)
	UpdateEmailService(input *InputUpdateEmail) (*model.EntityUsers, string)
}

type service struct {
	repository Repository
}

func NewAuthService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) ActivationService(input *InputActivation) (*model.EntityUsers, string) {
	users := model.EntityUsers{
		Email:  input.Email,
		Active: input.Active,
	}

	activationResult, activationError := s.repository.ActivationRepository(&users)

	return activationResult, activationError
}

func (s *service) ForgotService(input *InputForgot) (*model.EntityUsers, string) {

	users := model.EntityUsers{
		Email: input.Email,
	}

	resultRegister, errRegister := s.repository.ForgotRepository(&users)

	return resultRegister, errRegister
}

func (s *service) LoginService(input *InputLogin) (*model.EntityUsers, string) {

	user := model.EntityUsers{
		Email:    input.Email,
		Password: input.Password,
	}

	resultLogin, errLogin := s.repository.LoginRepository(&user)

	return resultLogin, errLogin
}

func (s *service) RegisterService(input *InputRegister) (*model.EntityUsers, string) {

	users := model.EntityUsers{
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: input.Password,
	}

	resultRegister, errRegister := s.repository.RegisterRepository(&users)

	return resultRegister, errRegister
}

func (s *service) ResendService(input *InputResend) (*model.EntityUsers, string) {

	users := model.EntityUsers{
		Email: input.Email,
	}

	resultRegister, errRegister := s.repository.ResendRepository(&users)

	return resultRegister, errRegister
}

func (s *service) ResetService(input *InputReset) (*model.EntityUsers, string) {

	users := model.EntityUsers{
		Email:    input.Email,
		Password: input.Password,
		Active:   input.Active,
	}

	resetResult, errResult := s.repository.ResetRepository(&users)

	return resetResult, errResult
}

func (s *service) UserService(email string) (*model.EntityUsers, string) {

	users := model.EntityUsers{
		Email: email,
	}

	resetResult, errResult := s.repository.UserRepository(&users)

	return resetResult, errResult
}

func (s *service) UpdateEmailService(input *InputUpdateEmail) (*model.EntityUsers, string) {

	// query to search engines

	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)

	go func() {
		c1 <- SearchGoogle("")
	}()
	go func() {
		c2 <- SearchYahoo("")
	}()
	go func() {
		c3 <- SearchBing("")
	}()

	result := ""
	for i := 0; i < 3; i++ {
		select {
		case msg1 := <-c1:
			result += msg1 + "_"
		case msg2 := <-c2:
			result += msg2 + "_"
		case msg3 := <-c3:
			result += msg3 + "_"
		}
	}
	fmt.Println(result)
	metaData := result[0 : len(result)-1]
	updateResult, updateError := s.repository.UpdateEmailRepository(input.Email, input.NewEmail, metaData)
	return updateResult, updateError
}

func SearchGoogle(msg string) string {
	time.Sleep(1 * time.Second)
	return "google"
}

func SearchYahoo(msg string) string {
	time.Sleep(1 * time.Second)
	return "yahoo"
}

func SearchBing(msg string) string {
	time.Sleep(2 * time.Second)
	return "bing"
}
