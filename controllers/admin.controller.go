package controllers

import (
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

func SignupAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		admin, err := CreateUser(user, false)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(
			utils.UserData{
				ID:   admin.ID,
				Role: admin.Role,
			},
		)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "Admin created successfully",
			"success": true,
			"token":   token,
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
