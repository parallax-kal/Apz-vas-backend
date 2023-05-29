package middlewares

import (
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

type Password struct {
	Password string `json:"password"`
}

func AdminAccountSettingsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		admin := c.MustGet("admin").(models.Admin)
		// get passwordBody from body
		var passwordBody Password

		if err := c.ShouldBindJSON(&passwordBody); err != nil {
			c.JSON(400, gin.H{
				"error":   "Password is required",
				"success": false,
			})
			c.Abort()
			return
		}
		err := utils.ComparePassword(passwordBody.Password, admin.Password)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "Incorrect Password",
				"success": false,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func OrganizationAccountSettingsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		organization := c.MustGet("organization").(models.Organization)
		var passwordBody Password

		if err := c.ShouldBindJSON(&passwordBody); err != nil {
			c.JSON(400, gin.H{
				"error":   "Password is required",
				"success": false,
			})
			c.Abort()
			return
		}
		err := utils.ComparePassword(passwordBody.Password, organization.Password)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "Incorrect Password",
				"success": false,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
