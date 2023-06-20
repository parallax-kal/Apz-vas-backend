package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
)

func CheckIfPaymentCanBeDone() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet = c.MustGet("wallet_data").(map[string]interface{})
		var organization = c.MustGet("organization_data").(models.Organization)
		var service = c.MustGet("service_data").(models.VASService)
		// https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/transfers

		var requestBody = make(map[string]interface{})
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// Handle the error
			c.JSON(400, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			// Handle the error
			c.JSON(400, gin.H{"error": "Failed to bind request body to struct"})
			c.Abort()
			return
		}
		var walletBalance = wallet["currentBalance"].(float64)
		var requestBodyAmount = requestBody["amount"].(float64)
		if walletBalance < requestBodyAmount {
			c.JSON(400, gin.H{
				"error":   "Insufficient Funds on Your wallet",
				"success": false,
			})
			c.Abort()
			return
		}

		var transferBody = make(map[string]interface{})

		transferBody["ignoreLimits"] = false
		transferBody["location"] = c.ClientIP()
		transferBody["replyPolicy"] = "WHEN_COMPLETE"
		transferBody["onlyCheck"] = true
		transferBody["fromWalletId"] = wallet["walletId"].(float64)
		transferBody["toWalletId"] = uint32(utils.ConvertStringToInt(os.Getenv("APZ_VAS_WALLET_ID")))
		transferBody["description"] = "Payment for " + service.Name
		transferBody["externalId"] = organization.ID
		transferBody["externalUniqueId"] = uuid.New()
		transferBody["amount"] = (requestBodyAmount * service.Rebate / 100) + requestBodyAmount
		transferBody["check"] = "BOTH"

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		var response, ukesheResponseError = UkhesheClient.Post("/wallets/transfers", transferBody)

		if ukesheResponseError != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   ukesheResponseError.Error(),
			})
			c.Abort()
			return
		}

		if response.Status != 204 {

			var transferResponse []map[string]interface{}

			json.Unmarshal(response.Data, &transferResponse)

			fmt.Println(transferResponse)

			c.JSON(500, gin.H{
				"success": false,
				"error":   transferResponse[0]["description"],
			})
			c.Abort()
			return
		}
		c.Set("request_body", requestBody)
		c.Next()
	}
}

func WalletMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		organization := c.MustGet("organization_data").(models.Organization)
		var wallet models.Wallet

		if err := configs.DB.Where("organization_id = ?", organization.ID).First(&wallet).Error; err != nil {
			c.JSON(403, gin.H{
				"error":   "You don't have a Wallet. Create one It is easy",
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
				"error":   "An error occured. Please try again or contact admin.",
				"success": false,
			})
			c.Abort()
			return
		}

		if response.Status != 200 {
			c.JSON(500, gin.H{
				"error":   "An error occured. Please try again or contact admin.",
				"success": true,
			})
			c.Abort()
			return
		}

		var wallet_body map[string]interface{}

		json.Unmarshal((response.Data), &wallet_body)

		if wallet_body["status"].(string) != "ACTIVE" {
			c.JSON(403, gin.H{
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
