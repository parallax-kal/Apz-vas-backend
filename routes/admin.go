package routes

import (
	"apz-vas/controllers"
	"github.com/gin-gonic/gin"
)

// initialize organization routes
func InitializeAdminRoutes(router *gin.RouterGroup) {
	router.POST("/signup", controllers.SignupAdmin())
}
