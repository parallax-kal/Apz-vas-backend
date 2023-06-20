package controllers

import (
	// "apz-vas/configs"
	// "apz-vas/models"
	// "apz-vas/utils"

	"github.com/gin-gonic/gin"
)

func SignupAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// var user models.User
		// if err := ctx.ShouldBindJSON(&user); err != nil {
		// 	ctx.JSON(400, gin.H{
		// 		"error":   err.Error(),
		// 		"success": false,
		// 	})
		// 	return
		// }
		// user.Role = "Admin"
		// newUser, err := CreateUser(user)

		// if err != nil {
		// 	ctx.JSON(400, gin.H{
		// 		"error":   err.Error(),
		// 		"success": false,
		// 	})
		// 	return
		// }

		// var admin models.Admin

		// if err := ctx.ShouldBindJSON(&admin); err != nil {
		// 	ctx.JSON(400, gin.H{
		// 		"error":   err.Error(),
		// 		"success": false,
		// 	})
		// 	return
		// }

		// admin.UserId = newUser.ID
		// admin.Role = "Admin"

		// if err := configs.DB.Create(&admin).Error; err != nil {
		// 	ctx.JSON(400, gin.H{
		// 		"error":   err.Error(),
		// 		"success": false,
		// 	})
		// 	return
		// }

		// token, err := utils.GenerateToken(
		// 	utils.UserData{
		// 		ID: newUser.ID,
		// 		Role: newUser.Role,
		// 	},
		// )
		// if err != nil {
		// 	ctx.JSON(500, gin.H{
		// 		"error":   "Something went wrong",
		// 		"success": false,
		// 	})
		// 	return
		// }
		// ctx.JSON(201, gin.H{
		// 	"message": "Admin created successfully",
		// 	"success": true,
		// 	"token":   token,
		// })

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
