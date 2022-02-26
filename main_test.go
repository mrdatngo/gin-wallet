package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	util "github.com/mrdatngo/gin-wallet/utils"
	"github.com/restuwahyu13/go-supertest/supertest"
	"syreclabs.com/go/faker"
)

var router = SetupRouter()
var accessToken string
var studentId interface{}

func TestLoginHandler(t *testing.T) {

	Convey("Auth Login Handler Group", t, func() {

		Convey("Login User Account Is Not Registered", func() {
			payload := gin.H{
				"email":    "anto13@zetmail.com",
				"password": "qwerty12",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/login")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusNotFound, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "User account is not registered", response.Message)
			})
		})

		Convey("Login Failed User Account Is Not Active", func() {
			payload := gin.H{
				"email":    "carmelo_marquardt@weissnat.info",
				"password": "testing13",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/login")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusForbidden, rr.Code)
				assert.Equal(t, http.MethodPost, response.Method)
				assert.Equal(t, "User account is not active", response.Message)
			})
		})

		Convey("Login Error Username Or Password Is Wrong", func() {
			payload := gin.H{
				"email":    "eduardo.wehner@greenholtadams.net",
				"password": "testing",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/login")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusForbidden, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Username or password is wrong", response.Message)
			})
		})

		Convey("Login Success", func() {
			payload := gin.H{
				"email":    "eduardo.wehner@greenholtadams.net",
				"password": "qwerty12345",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/login")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Login successfully", response.Message)

				var token map[string]interface{}
				encoded := util.Strigify(response.Data)
				_ = json.Unmarshal(encoded, &token)

				accessToken = token["accessToken"].(string)
			})
		})
	})
}

func TestRegisterHandler(t *testing.T) {

	Convey("Auth Register Handler Group", t, func() {

		Convey("Register New Account", func() {
			payload := gin.H{
				"fullname": faker.Internet().Email(),
				"email":    faker.Internet().Email(),
				"password": "testing13",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/register")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusCreated, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Register new account successfully", response.Message)
			})
		})

	})
}

func TestForgotHandler(t *testing.T) {

	Convey("Auth Forgot Password Handler Group", t, func() {

		Convey("Forgot Password If Email Not Exist", func() {

			payload := gin.H{
				"email": "santosi131@zetmail.com",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/forgot-password")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusNotFound, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Email is not never registered", response.Message)
			})
		})

		Convey("Forgot Password If Account Is Not Active", func() {

			payload := gin.H{
				"email": "santoso13@zetmail.com",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/forgot-password")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusForbidden, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "User account is not active", response.Message)
			})
		})

		Convey("Forgot Password To Get New Password", func() {

			payload := gin.H{
				"email": "samsul1@zetmail.com",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/forgot-password")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Forgot password successfully", response.Message)
			})
		})

	})
}

func TestResendHandler(t *testing.T) {

	Convey("Auth Resend Token Handler Group", t, func() {

		Convey("Resend New Token If Email Not Exist", func() {

			payload := gin.H{
				"email": "santosi131@zetmail.com",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/resend-token")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusNotFound, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Email is not never registered", response.Message)
			})
		})

		Convey("Resend Token If Account Is Active", func() {

			payload := gin.H{
				"email": "restuwahyu13@zetmail.com",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/resend-token")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusForbidden, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "User account hash been active", response.Message)
			})
		})

		Convey("Forgot Password To Get New Password", func() {

			payload := gin.H{
				"email": "santoso13@zetmail.com",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/resend-token")
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Resend new activation token successfully", response.Message)
			})
		})

	})
}

func TestResetHandler(t *testing.T) {

	Convey("Auth Reset Password Handler Group", t, func() {

		Convey("Reset Old Password To New Password", func() {
			payload := gin.H{
				"email":     "eduardo.wehner@greenholtadams.net",
				"password":  "qwerty12345",
				"cpassword": "qwerty12345",
			}

			test := supertest.NewSuperTest(router, t)

			test.Post("/api/v1/change-password/" + accessToken)
			test.Send(payload)
			test.Set("Content-Type", "application/json")
			test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

				response := util.Parse(rr.Body.Bytes())
				t.Log(response)

				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, "Change new password successfully", response.Message)
			})
		})

	})
}
