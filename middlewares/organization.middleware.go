package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

func OrganizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.MustGet("data").(utils.Data)
		var organization models.Organization
		if err := configs.DB.Where("ID = ?", data.ID).First(&organization).Error; err != nil {
			c.JSON(403, gin.H{
				"error":   "Unkown Organization",
				"success": false,
			})
			c.Abort()
			return
		}
		if organization.Status != "Active" {
			c.JSON(403, gin.H{
				"error":   "Inactive Organization",
				"success": false,
			})
		}
		c.Set("organization", organization)
		c.Next()
	}
}
