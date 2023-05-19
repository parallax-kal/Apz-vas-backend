package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		data := c.MustGet("data").(utils.Data)
		var admin models.Admin
		// search by ID
		if err := configs.DB.Where("ID = ?", data.ID).First(&admin).Error; err != nil {
			c.JSON(403, gin.H{
				"error":   "Unkown Admin",
				"success": false,
			})
			c.Abort()
			return
		}

		if admin.Status != "Active" {
			c.JSON(403, gin.H{
				"error":   "Inactive Admin",
				"success": false,
			})
			c.Abort()
			return
		}

		c.Set("admin", admin)

		c.Next()
	}
}
