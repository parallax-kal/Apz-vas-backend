package mobile

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAirtimeVendors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"airtime_vendors": []map[string]interface{}{
				{

					"id":   "cellc",
					"name": "Cell C",
				},
				{
					"id":   "mtn",
					"name": "MTN",
				},
				{
					"id":   "telkom",
					"name": "Telkom Mobile",
				},
				{
					"id":   "vodacom",
					"name": "Vodacom",
				},
			},
			"message": "Airtime Vendors Retrieved successfully",
			"success": true,
		})
	}
}

func BuyAirtime() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get mobile number, vendorId(cellc, mtn, telkom, vodacom), amount, deviceId, vendorId(used to identity a vendor in the blue label api) from the body

		var mobileNumber = c.PostForm("mobile_number")
		var vendorId = c.PostForm("vendor_id")
		var amount = c.PostForm("amount")
		var deviceId = c.PostForm("device_id")

		// get the organization id from the context
		var organization = c.MustGet("user_data").(*models.User)

		// get unique identifier for the request from client(this is not in the form submitted)

		if mobileNumber == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Mobile number is required",
			})
			return
		}

		if amount == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Amount is required",
			})
			return
		}

		var amountInt = utils.ConvertStringToInt(amount)

		if amountInt <= 0 {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Invalid amount",
			})
			return
		}

		if vendorId == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Vendor ID is required",
			})
			return
		}

		if deviceId == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Device ID is required",
			})
			return
		}

		// get the time of the transaction
		var transactionTime = time.Now().Format("2006-01-02T15:04:05") + "+02:00"
		var ipAddress = c.ClientIP()
		var payload = map[string]interface{}{
			"requestId":    ipAddress,
			"vendorId":     vendorId,
			"mobileNumber": mobileNumber,
			"amount":       amountInt,
			"vendMetaData": map[string]interface{}{
				"transactionRequestDateTime": transactionTime,
				// "transactionReference":       "0123456789",
				"vendorId": organization.ID,
				"deviceId": deviceId,
				// "consumerAccountNumber":      "012345",
			},
		}

		var blueLabelClient = configs.GetBlueLabelClient()

		var response, err = blueLabelClient.Post("/v2/trade/mobile/airtime/sales", payload)

		if err != nil {
			c.JSON(500, err)
			return
		}

		var responseBody map[string]interface{}

		json.Unmarshal(response.Data, &responseBody)

		if response.Status != 201 {
			c.JSON(response.Status, responseBody)
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Mobile Airtime bought successfully",
			"data":    responseBody,
		})

	}
}
