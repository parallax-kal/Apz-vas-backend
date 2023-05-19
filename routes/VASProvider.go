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
	// router.GET("/get-vas-providers", controllers.GetVASProviders())
}
