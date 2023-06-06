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

	router.GET("/get-vas-services", middlewares.AuthMiddleware(), controllers.GetVASServices())

	router.PUT("/update-vas-service",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.UpdateVasService(),
	)

	router.POST("/subscribe-vas-service", middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.SubScribeService(),
	)

	router.DELETE("/delete-vas-service",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.DeleteVasService(),
	)
}
