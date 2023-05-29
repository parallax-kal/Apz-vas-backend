package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

// initialize organization routes
func InitializeOrganizationRoutes(router *gin.RouterGroup) {
	router.POST("/signup", controllers.SignupOrganization())
	router.POST("/login", controllers.LoginOrganization())
	router.POST("/create-organization",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.CreateOrganization(),
	)
	router.GET("/get-organizations",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetOrganizations(),
	)
	router.GET("/subscribedServices",
		middlewares.OrganizationAPIMiddleware(),
		controllers.GetOrganizationSubScribedServices(),
	)
	router.POST("/subscribeService",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.SubScribeService(),
	)
	router.PUT("/account-settings",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.AdminAccountSettings(),
	)
}
