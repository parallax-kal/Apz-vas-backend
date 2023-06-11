package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		user := c.MustGet("user_data").(models.User)

		var admin models.Admin

		if err := configs.DB.Where("user_id = ?", user.ID).First(&admin).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"success": false,
			})
			c.Abort()
			return
		}

		// check if user's role is not admin or superadmin
		if admin.Role != "Admin" && admin.Role != "SuperAdmin" {
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
