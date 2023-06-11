package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"

	"github.com/gin-gonic/gin"
)

func OrganizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user_data").(models.User)

		var organization models.Organization

		if err := configs.DB.Where("user_id = ?", user.ID).First(&organization).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"success": false,
			})
			c.Abort()
			return
		}

		c.Set("organization_data", organization)

		c.Next()

	}
}
