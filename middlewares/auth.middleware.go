package middleware

import (
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract the token from the Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }

        // Parse and validate the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Provide the secret key or verification function for token validation
            // Replace "YOUR_SECRET_KEY" with your actual secret key or implement custom token verification logic
            return []byte("YOUR_SECRET_KEY"), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Check if the token is valid
        if !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Set the user ID or other relevant information in the context for further processing
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims"})
            c.Abort()
            return
        }

        userID := claims["user_id"].(string)
        c.Set("user_id", userID)

        c.Next()
    }
}
