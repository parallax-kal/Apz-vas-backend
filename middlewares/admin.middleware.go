package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if the user has the admin role
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
