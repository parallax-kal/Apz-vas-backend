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
	router.GET("/subscribedServices",
		middlewares.OrganizationAPIMiddleware(),
		controllers.GetOrganizationSubScribedServices(),
	)
	router.POST("/subscribeService",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.SubScribeService(),
	)
}
