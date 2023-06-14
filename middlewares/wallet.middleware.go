package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func PayForServiceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet = c.MustGet("wallet_data").(map[string]interface{})
		var organization = c.MustGet("organization_data").(models.Organization)
		var service = c.MustGet("service_data").(models.VASService)
		// https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/transfers
		var transfer models.Transaction
		if err := c.ShouldBindJSON(&transfer); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			c.Abort()
			return
		}
		if wallet["currentBalance"].(float64) < float64(transfer.Amount) {
			c.JSON(400, gin.H{
				"error":   "Insufficient Funds",
				"success": false,
			})
			c.Abort()
			return
		}
		var serviceData = make(map[string]interface{})
		
		if err := c.ShouldBindJSON(&serviceData); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			c.Abort()
			return
		}
		transfer.ServiceData = serviceData
		transfer.Amount = (uint64(service.Rebate) * uint64(transfer.Amount) / 100) + uint64(transfer.Amount)
		transfer.Location = c.ClientIP()
		transfer.ServiceId = service.ID
		transfer.Rebate = service.Rebate
		transfer.Currency = wallet["currency"].(string)
		transfer.ExternalId = organization.ID
		transfer.OtherWalletId = uint32(utils.ConvertStringToInt(os.Getenv("APZ_VAS_WALLET_ID")))
		transfer.UkhesheWalletId = wallet["walletId"].(float64)
		transfer.Description = "Payment for " + service.Name
		transfer.WalletId = utils.ConvertStringToUUID(wallet["externalUniqueId"].(string))

		if err := configs.DB.Create(&transfer).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			c.Abort()
			return
		}

		var transferBody = make(map[string]interface{})

		transferBody["ignoreLimits"] = true
		transferBody["location"] = transfer.Location
		transferBody["replyPolicy"] = "WHEN_COMPLETE"
		transferBody["onlyCheck"] = false
		transferBody["fromWalletId"] = transfer.UkhesheWalletId
		transferBody["toWalletId"] = transfer.OtherWalletId
		transferBody["description"] = transfer.Description
		transferBody["desccheckription"] = "BOTH"
		transferBody["externalId"] = transfer.ExternalId
		transferBody["externalUniqueId"] = transfer.ID

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		var response, ukesheResponseError = UkhesheClient.Post("/wallets/"+utils.ConvertIntToString(int(wallet["walletId"].(float64)))+"/transfers", transferBody)

		if ukesheResponseError != nil {
			configs.DB.Where("id = ?", transfer.ID).Delete(&transfer)
			c.JSON(500, gin.H{
				"success": false,
				"error":   ukesheResponseError.Error(),
			})
			c.Abort()
			return
		}

		if response.Status != 200 {
			configs.DB.Where("id = ?", transfer.ID).Delete(&transfer)

			var transferResponse []map[string]interface{}

			json.Unmarshal(response.Data, &transferResponse)

			c.JSON(500, gin.H{
				"success": false,
				"error":   transferResponse[0]["description"],
			})
			c.Abort()
			return
		}

		var transferResponse map[string]interface{}

		json.Unmarshal(response.Data, &transferResponse)
		// if transferResponse["status"].(string) != "COMPLETED" {
		// 	configs.DB.Where("id = ?", transfer.ID).Delete(&transfer)
		// 	c.JSON(500, gin.H{
		// 		"success": false,
		// 		"error":   "Something went wrong",
		// 	})
		// 	return
		// }

		c.Set("transfer_data", transfer)

		c.Next()
	}
}

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
