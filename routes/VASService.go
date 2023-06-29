package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeVASServiceRoutes(router *gin.RouterGroup) {

	router.GET("/get-vas-services",
		controllers.GetVasServices(),
	)

	router.GET("/get-organization-vas-services",
		middlewares.OrganizationAPIMiddleware(),
		controllers.GetOrganizationVASServices(),
	)

	router.GET("/get-admin-vas-services",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.GetAdminVASServices(),
	)

	router.PUT("/update-vas-service",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.UpdateVasService(),
	)

	router.DELETE("/delete-vas-service",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.DeleteVasService(),
	)

	router.PUT("/:operation",
		middlewares.AuthMiddleware(),
		middlewares.OrganizationMiddleware(),
		middlewares.VASServiceMiddleware(),
		controllers.OperationOnService(),
	)

	router.GET("/get-vas-service",
		middlewares.AuthMiddleware(),
		middlewares.VASServiceMiddleware(),
		controllers.GetVasServiceData(),
	)

	router.GET("/get-vas-service-transaction",
		middlewares.OrganizationAPIMiddleware(),
		middlewares.VASServiceMiddleware(),
		controllers.GetVasServiceTransactionHistory(),
	)

	router.GET("/get-admin-vas-service-transaction",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		middlewares.VASServiceMiddleware(),
		controllers.GetVasServiceAdminTransactionHistory(),
	)

}
