package routes

import (
	"apz-vas/controllers"
	// "apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

// initialize organization routes
func InitializeAdminRoutes(router *gin.RouterGroup) {
	router.POST("/signup", controllers.SignupAdmin())
	router.POST("/login", controllers.LoginAdmin())
}
