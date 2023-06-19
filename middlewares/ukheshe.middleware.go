package middlewares

import "github.com/gin-gonic/gin"

func CheckUkhesheClient() gin.HandlerFunc{
	return func(c*gin.Context) {
		// Check if the client is a valid Ukheshe client
		// If not, return an error
		// If yes, continue
		c.Next()
	}
}