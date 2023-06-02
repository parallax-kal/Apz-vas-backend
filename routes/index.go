package routes

import (
	"github.com/gin-gonic/gin"
)

// initialize all routes
func InitializeRoutes(router *gin.Engine) {
	organizationRoutes := router.Group("/organizations")
	adminRoutes := router.Group("/admins")
	vasProviderRoutes := router.Group("/vas-providers")
	vasServiceRoutes := router.Group("/vas-services")
	userRoutes := router.Group("/users")
	customerRoutes := router.Group("/customers")
	{
		InitializeAdminRoutes(adminRoutes)
		InitializeOrganizationRoutes(organizationRoutes)
		InitializeVASProviderRoutes(vasProviderRoutes)
		InitializeVASServiceRoutes(vasServiceRoutes)
		InitializeUserRouters(userRoutes)
		InitializeCustomerRoutes(customerRoutes)
	}
}
