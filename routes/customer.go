package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeCustomerRoutes(router *gin.RouterGroup) {
	router.GET("/get-customers",
		middlewares.OrganizationAPIMiddleware(),
		controllers.GetCustomers(),
	)
	router.POST("/create-customer",
		middlewares.OrganizationAPIMiddleware(),
		controllers.CreateCustomer(),
	)
}
