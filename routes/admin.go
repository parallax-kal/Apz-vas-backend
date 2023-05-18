package routes

import (
	"github.com/gin-gonic/gin"
	"apz-vas/middlewares"
	"apz-vas/controllers"
)

// initialize organization routes
func InitializeAdminRoutes(router *gin.RouterGroup) {
	router.POST("/create-organization", middlewares.AdminMiddleware(), controllers.CreateOrganization())
}