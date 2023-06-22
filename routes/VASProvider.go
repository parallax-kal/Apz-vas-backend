package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeVASProviderRoutes(router *gin.RouterGroup) {

	router.GET("/get-vas-providers",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetVasProviders(),
	)
	
	router.GET("/get-vas-provider-services",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetProviderServices(),
	)
	router.PUT("/update-vas-provider",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.UpdateVasProvider(),
	)
	router.PUT("/update-vas-provider-service",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.UpdateProviderService(),
	)

}
