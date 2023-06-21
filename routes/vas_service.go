package routes

import (
	"apz-vas/controllers/mobile"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeVasServiceRoutes(router *gin.RouterGroup) {
	router.GET("/airtime/get-vendors",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.NickNameService(),
		middlewares.ServiceProviderMiddleware(),
		mobile.GetAirtimeVendors(),
	)

	router.POST("/airtime/buy-airtime",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.NickNameService(),
		middlewares.CheckSubscription(),
		middlewares.ServiceProviderMiddleware(),
		middlewares.WalletMiddleware(),
		middlewares.CheckIfPaymentCanBeDone(),
		mobile.BuyAirtime(),
	)

	router.GET("/bundle/get-categories",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.NickNameService(),
		middlewares.ServiceProviderMiddleware(),
		mobile.GetMobileBundleCategories(),
	)

	router.GET("/bundle/get-products-by-category",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.NickNameService(),
		middlewares.ServiceProviderMiddleware(),
		mobile.GetMobileBundleProductsByCategory(),
	)

	router.POST("/bundle/buy-bundle",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.NickNameService(),
		middlewares.CheckSubscription(),
		middlewares.ServiceProviderMiddleware(),
		middlewares.WalletMiddleware(),
		middlewares.CheckIfPaymentCanBeDone(),
		mobile.BuyMobileBundle(),
	)

}
