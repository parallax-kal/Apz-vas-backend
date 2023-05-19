package routes

import (
	"apz-vas/controllers"
	"github.com/gin-gonic/gin"
)

// initialize organization routes
func InitializeOrganizationRoutes(router *gin.RouterGroup) {
	router.POST("/signup", controllers.SignupOrganization())
	router.POST("/login", controllers.LoginOrganization())
}
