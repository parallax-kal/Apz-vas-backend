package routes

import (
	"apz-vas/controllers"
	"github.com/gin-gonic/gin"
)

func InitializeVASProviderRoutes(router *gin.RouterGroup) {
	router.POST("/create-vas-provider", controllers.CreateVASProvider())
	// router.GET("/get-vas-providers", controllers.GetVASProviders())
}