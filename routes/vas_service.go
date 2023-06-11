package routes

import (
	"apz-vas/controllers/mobile"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeMobileRoutes(router *gin.RouterGroup) {
	router.GET("/airtime/get-vendors", middlewares.Check(), mobile.GetAirtimeVendors())

	router.POST("/airtime/buy-airtime", middlewares.Check(), middlewares.OrganizationAPIMiddleware(), mobile.BuyAirtime())

	router.GET("/bundle/categories", middlewares.Check(), mobile.GetMobileBundleCategories())

	router.GET("/bundle/get-products-by-category", mobile.GetMobileBundleProductsByCategory())

	router.POST("/bundle/buy-bundle", middlewares.OrganizationAPIMiddleware(), middlewares.Check(), mobile.BuyMobileBundle())

}
