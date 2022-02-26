package route

import (
	"github.com/gin-gonic/gin"
	auth_controllers "github.com/mrdatngo/gin-wallet/controllers/auth-controllers"
	auth_handlers "github.com/mrdatngo/gin-wallet/handlers/auth-handlers"
	middleware "github.com/mrdatngo/gin-wallet/middlewares"
	"gorm.io/gorm"
)

func InitAuthRoutes(db *gorm.DB, route *gin.Engine) {

	/**
	@description All Handler Auth
	*/
	authRepository := auth_controllers.NewAuthRepository(db)
	authService := auth_controllers.NewAuthService(authRepository)
	authHandler := auth_handlers.NewAuthHandler(authService)

	/**
	@description All Auth Route
	*/
	groupRoute := route.Group("/api/v1")
	groupRoute.POST("/register", authHandler.RegisterHandler)
	groupRoute.POST("/login", authHandler.LoginHandler)
	groupRoute.POST("/activation/:token", authHandler.ActivationHandler)
	groupRoute.POST("/resend-token", authHandler.ResendHandler)
	groupRoute.POST("/forgot-password", authHandler.ForgotHandler)
	groupRoute.POST("/change-password/:token", authHandler.ResetHandler)

	groupRoute.Use(middleware.Auth())
	groupRoute.POST("/update-email", authHandler.UpdateEmailHandler)
}
