package routes

import (
	"github.com/gin-gonic/gin"
)

// initialize all routes
func InitializeRoutes(router *gin.Engine) {
	organizationRoutes := router.Group("/organization")
	adminRoutes := router.Group("/admin")
	vasProviderRoutes := router.Group("/vas-providers")
	vasServiceRoutes := router.Group("/vas-services")
	userRoutes := router.Group("/users")
	{
		InitializeAdminRoutes(adminRoutes)
		InitializeOrganizationRoutes(organizationRoutes)
		InitializeVASProviderRoutes(vasProviderRoutes)
		InitializeVASServiceRoutes(vasServiceRoutes)
		InitializeUserRouters(userRoutes)
	}
}
