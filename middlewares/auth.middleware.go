package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

type ReturnedUser struct {
	// Name
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if err := configs.DB.Select("status, role, name, id, api_key", "email").Where("id = ?", userData.ID).First(&user).Error; err != nil {
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
		user.Status = ""
		c.Set("user_data", user)
		c.Next()
	}
}
