package controllers

import (
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
	"os"
)

func SignupAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user utils.UserEmailedData

		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{
				"error":   "Bad Request",
				"success": false,
			})
			return
		}

		user.Role = "Admin"

		if err := ValidateUser(user); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		token, err := utils.GenerateTokenFromUserData(user)

		if err != nil {
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var link = os.Getenv("FRONTEND_URL") + "/signup/continue?token=" + token

		if err := utils.SendMail(user.Email, "Admin Signup", "Please click on the link below to continue with your registration: "+link); err != nil {
			ctx.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Email sent successfully",
			"success": true,
		})

	}
}

func GetAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func DeleteAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
