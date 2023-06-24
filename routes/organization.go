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
		middlewares.AdminMiddleware(),
		controllers.CreateOrganization(),
	)

	router.POST("/signup",
		middlewares.PartialAuthMiddleware(),
		controllers.SignupOrganizationContinue(),
	)

	router.GET("/get-your-organization-data",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		controllers.GetYourOrganizationData(),
	)

	router.DELETE("/delete-organization",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.DeleteOrganization(),
	)
	router.PUT("/update-organization",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.UpdateOrganization(),
	)
	router.PUT("/settings",
		middlewares.PasswordChecker(),
		middlewares.OrganizationMiddleware(),
		controllers.UpdateOrganization(),
	)
	router.GET("/get-organizations",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetOrganizations(),
	)

}
