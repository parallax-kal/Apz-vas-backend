package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeVASServiceRoutes(router *gin.RouterGroup) {
	router.POST("/create-vas-service",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.CreateVasService(),
	)
	// router.GET("/get-vas-services", controllers.GetVASServices())
	// router.GET("/get-vas-service/:id", controllers.GetVASService())
	// router.PUT("/update-vas-service/:id", controllers.UpdateVASService())
	// router.DELETE("/delete-vas-service/:id", controllers.DeleteVASService())
}
