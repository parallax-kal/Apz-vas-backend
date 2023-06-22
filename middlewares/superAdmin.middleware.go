package middlewares

import (
	"apz-vas/models"
	"github.com/gin-gonic/gin"
)

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user_data").(models.User)
		if user.Role != "SuperAdmin" {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"success": false,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
