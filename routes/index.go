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

	mobileRoutes := router.Group("/vas-service/mobile")

	{
		InitializeOrganizationRoutes(organizationRoutes)
		InitializeUserRouters(userRoutes)
		InitializeVASServiceRoutes(vasServiceRoutes)
		InitializeVASProviderRoutes(vasProviderRoutes)
		InitializeAdminRoutes(adminRoutes)
		InitializeCustomerRoutes(customerRoutes)
		InitializeMobileRoutes(mobileRoutes)
	}
}
