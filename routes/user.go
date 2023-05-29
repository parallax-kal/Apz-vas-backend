package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeUserRouters(router *gin.RouterGroup) {
	router.POST("/login", controllers.LoginUser())
	router.PUT("/account-settings",
		middlewares.AccountSettingsMiddleware(),
		controllers.AccountSettings(),
	)
}
