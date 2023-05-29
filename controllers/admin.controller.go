package controllers

import (
	"apz-vas/configs"
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
		user.Role = "Admin"
		newUser, err := CreateUser(user)
		var admin models.Admin
		admin.UserId = newUser.ID
		// tx := configs.DB.Begin()
		if err :=  configs.DB.Create(&admin).Error; err != nil {
			// tx.Rollback()
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		// tx.Commit()
		token, err := utils.GenerateToken(
			utils.Data{
				ID: newUser.ID,
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

