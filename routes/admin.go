package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

// initialize organization routes
func InitializeAdminRoutes(router *gin.RouterGroup) {
	router.POST("/signup", controllers.SignupAdmin())
	router.GET("/get-admins",
		middlewares.SuperAdminMiddleware(),
		controllers.GetAdmins(),
	)
	router.PUT("/update-admin",
		middlewares.SuperAdminMiddleware(),
		controllers.UpdateAdmin(),
	)
	router.DELETE("/delete-admin",
		middlewares.SuperAdminMiddleware(),
		controllers.DeleteAdmin(),
	)
	router.POST("/create-admin",
		middlewares.SuperAdminMiddleware(),
		controllers.CreateAdmin(),
	)

}
