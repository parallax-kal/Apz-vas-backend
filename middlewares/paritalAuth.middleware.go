package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

func PartialAuthMiddleware() gin.HandlerFunc {
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

		if err := configs.DB.Select("status, role, name, id", "email").Where("id = ?", userData.ID).First(&user).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "Invalid token",
				"success": false,
			})
			c.Abort()
			return
		}

		user.Passwords = ""

		c.Set("user_data", user)
		c.Next()
	}
}
