package mobile

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var categories = []string{
	"data",
	"sms",
}

func GetMobileBundleCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success":    true,
			"categories": categories,
			"message":    "Mobile Bundle Categories Retrieved successfully",
		})
	}
}

func GetMobileBundleProductsByCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var category = c.Query("category")

		if category == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Category is required",
			})
			return
		}

		if !utils.Contains(categories, category) {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Invalid category provided use this (data, sms)",
			})
			return
		}

		var query = url.Values{
			"category": []string{category},
		}

		var response, err = configs.BlueLabelCleint.Get("/mobile/bundle/products", query)

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Error occurred while fetching mobile bundle products",
			})
			return
		}

		var responseBody []map[string]interface{}

		json.Unmarshal(response.Data, &responseBody)

		if response.Status != 200 {
			c.JSON(response.Status, responseBody)
			return
		}

		c.JSON(200, gin.H{
			"products": responseBody,
			"message":  "Mobile Bundle Products Retrieved successfully",
			"success":  true,
		})

	}
}

func BuyMobileBundle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody = c.MustGet("request_body").(map[string]interface{})
		mobileNumber, vendorId, deviceId, productId := requestBody["mobile_number"].(string), requestBody["vendor_id"].(string), requestBody["device_id"].(string), requestBody["product_id"].(string)
		var organization = c.MustGet("organization_data").(models.Organization)

		if mobileNumber == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Mobile number is required",
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

		if productId == "" {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Product ID is required",
			})
		}

		var transactionTime = time.Now().Format("2006-01-02T15:04:05") + "+02:00"
		var ipAddress = c.ClientIP()
		var payload = map[string]interface{}{
			"requestId":    ipAddress,
			"vendorId":     vendorId,
			"productId":    productId,
			"mobileNumber": mobileNumber,
			"vendMetaData": map[string]interface{}{
				"transactionRequestDateTime": transactionTime,
				// "transactionReference":       "0123456789",
				"vendorId": organization.ID,
				"deviceId": deviceId,
				// "consumerAccountNumber":      "012345",
			},
		}

		var response, err = configs.BlueLabelCleint.Post("/mobile/bundle/sales", payload)

		if err != nil {
			c.JSON(500, err)
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

		utils.PayForService(c)

		c.JSON(200, gin.H{
			"success": true,
			"message": "Mobile Bundle bought successfully.",
			"data":    responseBody,
		})

	}
}
