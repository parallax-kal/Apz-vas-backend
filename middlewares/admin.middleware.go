package middlewares

import (
	"apz-vas/models"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := c.MustGet("user_data").(models.User)

		// check if user's role is not admin or superadmin
		if user.Role != "Admin" && user.Role != "SuperAdmin" {
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
