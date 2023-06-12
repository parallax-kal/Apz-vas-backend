package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeWalletRoutes(router *gin.RouterGroup) {

	router.GET("/get-wallet-types",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.GetWalletTypes(),
	)

	router.GET("/get-wallet",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.GetWallet(),
	)

	router.POST("/create-wallet",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.CreateWallet(),
	)

	router.POST("/topup-wallet",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.TopUpWallet(),
	)

	router.GET("/get-transaction-history",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
	)

	router.POST("/withdraw-wallet",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.WithDrawWallet(),
	)

}
