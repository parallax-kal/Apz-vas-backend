package middlewares

import (
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{
				"error":   "Token is missing",
				"success": false,
			})
			c.Abort()
			return
		}
		data, error := utils.ExtractDataFromToken(tokenString)
		if error != nil {
			c.JSON(401, gin.H{
				"error":   error.Error(),
				"success": false,
			})
			return
		}

		c.Set("data", data)
		c.Next()
	}
}
