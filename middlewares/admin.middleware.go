package middlewares

import (
	"apz-vas/utils"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user has the admin role

		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		data, error := utils.ExtractDataFromToken(tokenString)
        if error != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        }

        fmt.Println(data)

		// Implement your logic to check the admin role based on the user role or other criteria
		isAdmin := true // Example logic, modify as per your requirements
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
