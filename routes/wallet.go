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
		middlewares.WalletMiddleware(),
		controllers.GetWallet(),
	)

	router.GET("/get-wallet-balances",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.GetWalletBalances(),
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
		middlewares.OrganizationAPIMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.GetTransactionHistory(),
	)

	router.GET("/get-transaction-history/topup",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.GetTopupHistory(),
	)
	router.GET("/get-transaction-history/withdraw",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.GetWithdrawHistory(),
	)

	router.POST("/withdraw-from-wallet",
		middlewares.PasswordChecker(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.WithDrawFromWallet(),
	)

	router.PUT("/update-wallet",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.UpdateWallet(),
	)

	router.GET("/get-withdrawal-fees",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.WalletMiddleware(),
		controllers.GetWithdrawalFees(),
	)

	router.POST("/topup-callabck",
		middlewares.CheckUkhesheClient(),
		controllers.UkhesheTopupCallBack(),
	)

	router.POST("/withdraw-callabck",
		middlewares.CheckUkhesheClient(),
		controllers.UkhesheWithdrawCallBack(),
	)

}
