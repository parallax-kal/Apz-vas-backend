package middlewares

import (
	"apz-vas/models"
	"github.com/gin-gonic/gin"
)

func OrganizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		organization := c.MustGet("user_data").(models.User)
		if organization.Role != "Organization" {
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
