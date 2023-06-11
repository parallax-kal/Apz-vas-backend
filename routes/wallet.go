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

}
