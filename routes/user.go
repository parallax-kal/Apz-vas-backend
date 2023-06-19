package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeUserRouters(router *gin.RouterGroup) {
	router.POST("/login", controllers.LoginUser())
	router.PUT("/account-settings",
		middlewares.AuthMiddleware(),
		middlewares.PasswordChecker(),
		controllers.AccountSettings(),
	)
	router.PUT("/change-password",
		middlewares.AuthMiddleware(),
		middlewares.PasswordChecker(),
		controllers.UpdatePassword(),
	)
	router.POST("/signup", controllers.SignupOrganization())
	router.POST("/google-login", controllers.GoogleLogin())
	router.POST("/google-signup", controllers.GoogleRegister())
	router.GET("/get-users",
		middlewares.AuthMiddleware(),
		middlewares.SuperAdminMiddleware(),
		controllers.GetUsers(),
	)
	router.GET("/me",
		middlewares.PartialAuthMiddleware(),
		controllers.GetUser(),
	)
}
