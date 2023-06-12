package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"

	"github.com/gin-gonic/gin"
)

func WalletMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var walletId = c.Query("walletId")

		if walletId == "" {
			c.JSON(404, gin.H{
				"error":   "No Wallet Provided",
				"success": false,
			})
			c.Abort()
			return
		}

		var wallet models.Wallet
		if err := configs.DB.Where("id = ?", walletId).First(&wallet).Error; err != nil {
			c.JSON(404, gin.H{
				"error":   "Unknown Wallet",
				"success": true,
			})
			c.Abort()
			return
		}
		c.Set("wallet_data", wallet)
		c.Next()
	}
}
