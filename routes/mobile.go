package routes

import (
	"apz-vas/controllers/mobile"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeMobileRoutes(router *gin.RouterGroup) {
	router.GET("/get-mobile-vendors", mobile.GetAirtimeVendors())

	router.POST("/buy-mobile-airtime", middlewares.OrganizationAPIMiddleware(), mobile.BuyAirtime())

	router.GET("/get-mobile-bundle-categories", mobile.GetMobileBundleCategories())

	router.GET("/get-mobile-bundle-products-by-category", mobile.GetMobileBundleProductsByCategory())

	router.POST("/buy-mobile-bundle", middlewares.OrganizationAPIMiddleware(), mobile.BuyMobileBundle())

}
