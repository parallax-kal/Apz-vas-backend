package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"

	"github.com/gin-gonic/gin"
)

func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user_data").(models.User)

		var superAdmin models.Admin

		if err := configs.DB.Where("user_id = ?", user.ID).First(&superAdmin).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"success": false,
			})
			c.Abort()
			return
		}

		if superAdmin.Role != "SuperAdmin" {
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
