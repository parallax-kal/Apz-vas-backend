package routes

import (
	"apz-vas/controllers"
	"apz-vas/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeUserRouters(router *gin.RouterGroup) {
	router.POST("/login", controllers.LoginUser())
	router.PUT("/account-settings",
		middlewares.PasswordChecker(),
		controllers.AccountSettings(),
	)
	router.PUT("/change-password",
		middlewares.PasswordChecker(),
		controllers.UpdatePassword(),
	)
	router.POST("/get-verification-link",
		middlewares.PasswordChecker(),
		controllers.GetVerificationLink(),
	)
	router.PUT("/change-email",
		controllers.ChangeEmail(),
	)
	router.POST("/signup", controllers.SignupOrganization())
	router.GET("/verify", controllers.VerifyUser())
	router.GET("/verify-reseting", controllers.VerifyBeforeResetingPassword())
	router.PUT("/reset-password", controllers.ResetPassword())
	router.POST("/forgot-password", controllers.ForgotPassword())
	router.POST("/google-login", controllers.GoogleLogin())
	router.POST("/google-signup", controllers.GoogleRegister())
	router.GET("/get-users",
		middlewares.AuthMiddleware(),
		middlewares.SuperAdminMiddleware(),
		controllers.GetUsers(),
	)
	router.POST("/validate-org-data",
		middlewares.AuthMiddleware(),
		middlewares.AdminMiddleware(),
		controllers.ValidateUserOrgData(),
	)
	router.GET("/me",
		middlewares.PartialAuthMiddleware(),
		controllers.GetUser(),
	)
}
