package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func WalletMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		organization := c.MustGet("organization_data").(models.Organization)
		var wallet models.Wallet

		if err := configs.DB.Where("organization_id = ?", organization.ID).First(&wallet).Error; err != nil {
			c.JSON(200, gin.H{
				"error":   "Organization don't have a Wallet",
				"wallet":  nil,
				"success": true,
			})
			c.Abort()
			return
		}

		// https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/{walletId}
		Ukheshe_Client := configs.MakeAuthenticatedRequest(true)
		response, err := Ukheshe_Client.Get("/wallets/" + utils.ConvertIntToString(int(wallet.Ukheshe_Id)))

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{
				"error":   "Something Went Wrong",
				"success": false,
			})
			c.Abort()
			return
		}

		if response.Status != 200 {
			c.JSON(404, gin.H{
				"error":   "Wallet Not Found",
				"success": true,
			})
			c.Abort()
			return
		}

		var wallet_body map[string]interface{}

		json.Unmarshal((response.Data), &wallet_body)

		if wallet_body["status"].(string) != "ACTIVE" {
			c.JSON(200, gin.H{
				"error":   "Your Wallet is not Active",
				"success": true,
			})
			c.Abort()
			return
		}

		c.Set("wallet_data", wallet_body)
		c.Next()
	}
}
