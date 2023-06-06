package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

// initialize organization routes
func InitializeOrganizationRoutes(router *gin.RouterGroup) {
	router.POST("/create-organization",
		middlewares.AuthMiddleware(),
		controllers.CreateOrganization(),
	)
	router.POST("/signup", controllers.SignupOrganization())
	router.DELETE("/delete-organization",
		middlewares.AdminMiddleware(),
		controllers.DeleteOrganization(),
	)
	router.PUT("/update-organization",
		middlewares.AdminMiddleware(),
		controllers.UpdateOrganization(),
	)
	router.GET("/get-organizations",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetOrganizations(),
	)

}
