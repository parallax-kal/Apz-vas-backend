package mobile

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

func GetAirtimeVendors() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response, err = configs.BlueLabelCleint.Get("/mobile/airtime/products")
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "Something went wrong, try again or contact admin.",
				"success": false,
			})
			return
		}
		var responseBody []map[string]interface{}
		if err := json.Unmarshal(response.Data, &responseBody); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"airtime_vendors": responseBody,
			"message":         "Airtime Vendors Retrieved successfully",
			"success":         true,
		})
	}
}

func BuyAirtime() gin.HandlerFunc {
	return func(c *gin.Context) {

		var requestBody = c.MustGet("request_body").(map[string]interface{})

		mobileNumber, amount, deviceId, vendorId := requestBody["mobile_number"].(string), requestBody["amount"].(float64), requestBody["device_id"].(string), requestBody["vendor_id"].(string)
		var organization = c.MustGet("organization_data").(models.Organization)

		if mobileNumber == "" {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Mobile number is required",
			})
			return
		}

		if vendorId == "" {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Vendor ID is required",
			})
			return
		}

		if deviceId == "" {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Device ID is required",
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
			"amount":       amount,
			"vendMetaData": map[string]interface{}{
				"transactionRequestDateTime": transactionTime,
				// "transactionReference":       "0123456789",
				"vendorId": organization.ID,
				"deviceId": deviceId,
				// "consumerAccountNumber":      "012345",
			},
		}

		var response, err = configs.BlueLabelCleint.Post("/mobile/airtime/sales", payload)

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var responseBody map[string]interface{}

		json.Unmarshal(response.Data, &responseBody)

		if response.Status != 201 {
			c.JSON(response.Status, gin.H{
				"success": false,
				"error":   responseBody["message"],
			})
			return
		}
		paymentError := utils.PayForService(c)
		if paymentError != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   paymentError.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"success": true,
			"message": "Mobile Airtime bought successfully",
		})

	}
}
