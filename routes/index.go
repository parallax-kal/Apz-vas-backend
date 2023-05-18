package routes

import (
	"github.com/gin-gonic/gin"

)

// initialize all routes
func InitializeRoutes(router *gin.Engine) {
	organizationRoutes := router.Group("/organization")
	adminRoutes := router.Group("/admin")
	{
		InitializeAdminRoutes(adminRoutes)
		InitializeOrganizationRoutes(organizationRoutes)
	}
}