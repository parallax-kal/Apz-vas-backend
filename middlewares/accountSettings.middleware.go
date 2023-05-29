package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

type Password struct {
	OldPassword string `json:"password"`
}


func AccountSettingsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var user models.User
		var passwordBody Password
		if err := c.ShouldBindJSON(&passwordBody); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			c.Abort()
			return
		}
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{
				"error":   "Token is missing",
				"success": false,
			})
			c.Abort()
			return
		}
		userData, error := utils.ExtractDataFromToken(tokenString)
		if error != nil {
			c.JSON(401, gin.H{
				"error":   error.Error(),
				"success": false,
			})
			c.Abort()
			return
		}
		var user models.User		
		
		if err := configs.DB.Where("ID = ?", userData.ID).First(&user).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "User not found",
				"success": false,
			})
			c.Abort()
			return
		}

		if user.Status != "Active" {
			c.JSON(401, gin.H{
				"error":   "User is inactive",
				"success": false,
			})
			c.Abort()
			return
		}

		if err := utils.ComparePassword(passwordBody.OldPassword, user.Password); err != nil {
			c.JSON(401, gin.H{
				"error":   "Password is incorrect",
				"success": false,
			})
			c.Abort()
			return
		}
		c.Set("user_data", user)
		c.Next()

	}
}