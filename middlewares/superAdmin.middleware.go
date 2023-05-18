package middlewares

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func SuperAdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if the user has the superadmin role
        // Implement your logic to check the superadmin role based on the user role or other criteria
        isSuperAdmin := true // Example logic, modify as per your requirements
        if !isSuperAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}
