package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func OrganizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user_data").(models.User)

		var organization models.Organization

		if err := configs.DB.Where("user_id = ?", user.ID).First(&organization).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "Organization not found",
				"success": false,
			})
			c.Abort()
			return
		}


		Ukheshe_Client := configs.MakeAuthenticatedRequest(true)

		response, err := Ukheshe_Client.Get("/organisations/" + utils.ConvertIntToString(int(organization.Ukheshe_Id)))

		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Something Went Wrong",
				"success": false,
			})
			c.Abort()
			return
		}

		if response.Status != 200 {
			c.JSON(401, gin.H{
				"error":   "You account is not up to date, please contact admin",
				"success": false,
			})
			c.Abort()
			return
		}

		var response_data map[string]interface{}

		if err := json.Unmarshal(response.Data, &response_data); err != nil {
			c.JSON(500, gin.H{
				"error":   "Something Went Wrong",
				"success": false,
			})
			c.Abort()
			return
		}

		c.Set("organization_version", response_data["version"])

		c.Set("organization_data", organization)
		c.Next()

	}
}
