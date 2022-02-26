package route

import (
	"github.com/gin-gonic/gin"
	auth_controllers "github.com/mrdatngo/gin-wallet/controllers/auth-controllers"
	wallet_controllers "github.com/mrdatngo/gin-wallet/controllers/wallet-controllers"
	wallet_handlers "github.com/mrdatngo/gin-wallet/handlers/wallet-handlers"
	middleware "github.com/mrdatngo/gin-wallet/middlewares"
	"gorm.io/gorm"
)

func InitWalletRoutes(db *gorm.DB, route *gin.Engine) {

	/**
	@description All Handler Wallet
	*/

	authRepository := auth_controllers.NewAuthRepository(db)
	authService := auth_controllers.NewAuthService(authRepository)

	walletRepository := wallet_controllers.NewWalletRepository(db)
	walletService := wallet_controllers.NewWalletService(walletRepository)
	walletHandler := wallet_handlers.NewWalletHandler(walletService, authService)

	/**
	@description All Wallet Route
	*/
	groupRoute := route.Group("/api/v1").Use(middleware.Auth())
	groupRoute.GET("/wallets", walletHandler.ResultsWalletHandler)
	groupRoute.POST("/wallet", walletHandler.CreateWalletHandler)
	groupRoute.DELETE("/wallet/:id", walletHandler.DeleteWalletHandler)
}
