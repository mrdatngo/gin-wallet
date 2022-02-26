package route

import (
	"github.com/gin-gonic/gin"
	auth_controllers "github.com/mrdatngo/gin-wallet/controllers/auth-controllers"
	transaction_controllers "github.com/mrdatngo/gin-wallet/controllers/transaction-controllers"
	transaction_handlers "github.com/mrdatngo/gin-wallet/handlers/transaction-handlers"
	middleware "github.com/mrdatngo/gin-wallet/middlewares"
	"gorm.io/gorm"
)

func InitTransactionRoutes(db *gorm.DB, route *gin.Engine) {

	/**
	@description All Handler Student
	*/

	authRepository := auth_controllers.NewAuthRepository(db)
	authService := auth_controllers.NewAuthService(authRepository)

	transactionRepository := transaction_controllers.NewTransactionRepository(db)
	transactionService := transaction_controllers.NewTransactionService(transactionRepository)
	transactionHandler := transaction_handlers.NewTransactionHandler(transactionService, authService)

	/**
	@description All Student Route
	*/
	groupRoute := route.Group("/api/v1").Use(middleware.Auth())
	groupRoute.POST("/deposit", transactionHandler.DepositHandler)
	groupRoute.GET("/deposits", transactionHandler.DepositListHandler)
}
