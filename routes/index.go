package routes

import (
	"github.com/gin-gonic/gin"
)

// initialize all routes
func InitializeRoutes(router *gin.Engine) {
	organizationRoutes := router.Group("/organizations")
	adminRoutes := router.Group("/admins")
	vasProviderRoutes := router.Group("/vas-providers")
	vasServicesRoutes := router.Group("/vas-services")
	userRoutes := router.Group("/users")
	customerRoutes := router.Group("/customers")
	walletRoutes := router.Group("/wallets")
	vasServiceRoutes := router.Group("/vas-service")

	{
		InitializeOrganizationRoutes(organizationRoutes)
		InitializeUserRouters(userRoutes)
		InitializeVASServiceRoutes(vasServicesRoutes)
		InitializeVASProviderRoutes(vasProviderRoutes)
		InitializeAdminRoutes(adminRoutes)
		InitializeCustomerRoutes(customerRoutes)
		InitializeMobileRoutes(vasServiceRoutes)
		InitializeWalletRoutes(walletRoutes)
	}
}
