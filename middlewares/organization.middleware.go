package middlewares

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func OrganizationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implement your logic to check if the user belongs to the organization
        // You can access the organization information from the token or other sources

        // Example logic: Check if the user is part of the organization
        organizationID := c.MustGet("organization_id").(string)
        if organizationID != "user_organization" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}
