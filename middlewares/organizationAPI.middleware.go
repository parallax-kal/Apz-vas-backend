package middlewares

import (
	"apz-vas/configs"
	"apz-vas/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func OrganizationAPIMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apikey := c.Request.Header.Get("apz-vas-api-key")
		if apikey == "" {
			c.JSON(401, gin.H{
				"error":   "Api Key missing",
				"success": false,
			})
			c.Abort()
			return
		}
		APIKey := utils.ConvertStringToUUID(apikey)
		organization, err := utils.CheckApiKey(APIKey)
		if err != nil {
			c.JSON(401, gin.H{
				"error":   "Invalid Api Key",
				"success": false,
			})
			c.Abort()
			return
		}

		Ukheshe_Client := configs.MakeAuthenticatedRequest(true)

		response, err := Ukheshe_Client.Get("/organisations/" + utils.ConvertIntToString(int(organization.Ukheshe_Id)))

		if err != nil {
			c.JSON(500, gin.H{
				"error":   "An error occurred. Please try again or Contact admin",
				"success": false,
			})
			c.Abort()
			return
		}

		if response.Status != 200 {
			c.JSON(500, gin.H{
				"error":   "An error occurred. Please try again or Contact Admin.",
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

		c.Set("organization_data", organization)
		c.Next()

	}
}
