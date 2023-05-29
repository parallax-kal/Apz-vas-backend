package routes

import (
	"github.com/gin-gonic/gin"

)

// initialize all routes
func InitializeRoutes(router *gin.Engine) {
	organizationRoutes := router.Group("/organization")
	adminRoutes := router.Group("/admin")
	vasProviderRoutes := router.Group("/vas-provider")
	vasServiceRoutes := router.Group("/vas-service")
	userRoutes := router.Group("/user")
	{
		InitializeAdminRoutes(adminRoutes)
		InitializeOrganizationRoutes(organizationRoutes)
		InitializeVASProviderRoutes(vasProviderRoutes)
		InitializeVASServiceRoutes(vasServiceRoutes)
		InitializeUserRouters(userRoutes)
	}
}