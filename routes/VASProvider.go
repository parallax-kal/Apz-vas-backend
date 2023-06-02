package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeVASProviderRoutes(router *gin.RouterGroup) {
	router.POST("/create-vas-provider",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.CreateVASProvider(),
	)
	router.GET("/get-vas-providers",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetVasProviders(),
	)
	
	router.GET("/get-vas-provider-service",
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
	router.DELETE("/delete-vas-provider",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.DeleteVasProvider(),
	)
}
