package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetWalletTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		response, err := UkhesheClient.Get("/wallet-types")

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var responseBody []map[string]interface{}
		if response.Status != 200 {
			c.JSON(response.Status, gin.H{
				"success": false,
				"error":   responseBody[0]["description"],
			})
			return
		}

		if err := json.Unmarshal(response.Data, &responseBody); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// delete the wallet type of mode system
		for i := 0; i < len(responseBody); i++ {
			if responseBody[i]["mode"] == "SYSTEM" || responseBody[i]["mode"] == "PREPAID_CARD" {
				responseBody = append(responseBody[:i], responseBody[i+1:]...)
			}
			// delete mode, version, configuration
			delete(responseBody[i], "mode")
			delete(responseBody[i], "version")
			delete(responseBody[i], "configuration")
		}

		c.JSON(200, gin.H{
			"success":      true,
			"wallet_types": responseBody,
		})

	}

}

func CreateWallet() gin.HandlerFunc {
	return func(c *gin.Context) {

		var wallet models.Wallet

		if err := c.ShouldBindJSON(&wallet); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}
		organization := c.MustGet("organization_data").(models.Organization)

		wallet.OrganizationId = organization.ID

		if err := configs.DB.Create(&wallet).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		// fmt.Println(wallet)
		var walletBody = make(map[string]interface{})
		walletBody["cardType"] = wallet.CardType
		walletBody["description"] = wallet.Description
		walletBody["externalUniqueId"] = wallet.ID

		walletBody["name"] = wallet.Name

		// add status as uppercase to the wallet.status
		walletBody["status"] = strings.ToUpper(wallet.Status)
		walletBody["walletTypeId"] = wallet.WalletTypeID
		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		response, err := UkhesheClient.Post("/organisations/"+utils.ConvertIntToString(int(organization.Ukheshe_Id))+"/wallets", walletBody)

		if err != nil {
			// delete wallet
			configs.DB.Where("id = ?", wallet.ID).Delete(&wallet)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if response.Status != 200 {
			configs.DB.Where("id = ?", wallet.ID).Delete(&wallet)
			var responseBody []map[string]interface{}
			c.JSON(response.Status, gin.H{
				"success": false,
				"error":   responseBody[0]["description"],
			})
			return
		}

		var responseBody map[string]interface{}

		if err := json.Unmarshal(response.Data, &responseBody); err != nil {
			configs.DB.Where("id = ?", wallet.ID).Delete(&wallet)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// update wallet
		if err := configs.DB.Model(&models.Wallet{}).Where("id = ?", wallet.ID).Update("ukheshe_id", responseBody["walletId"]).Error; err != nil {
			configs.DB.Where("id = ?", wallet.ID).Delete(&wallet)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"wallet":  responseBody,
		})
	}

}

func GetWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet models.Wallet

		organization := c.MustGet("organization_data").(models.Organization)

		if err := configs.DB.Where("organization_id = ?", organization.ID).First(&wallet).Error; err != nil {
			if err.Error() == "record not found" {
				c.JSON(200, gin.H{
					"message": "No Wallet Created yet",
					"wallet":  nil,
					"success": true,
				})

			} else {

				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
			}

			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"wallet":  wallet,
		})
	}
}

type TopupCard struct {
	AccountType    string `json:"account_type"`
	Alias          string `json:"alisa"`
	CardHolderName string `json:"card_holder_name"`
	Cvv            string `json:"cvv"`
	Dob            string `json:"card_holder_dob"`
	Expiry         string `json:"expiry"`
	Pan            string `json:"pan"`
}

func WithDrawWallet() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func TopUpWallet() gin.HandlerFunc {

	return func(c *gin.Context) {
		var TopupData models.Topup

		if err := c.ShouldBindJSON(&TopupData); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		wallet := c.MustGet("wallet_data").(models.Wallet)

		TopupData.WalletId = wallet.ID

		if err := configs.DB.Create(&TopupData).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var topupBody = make(map[string]interface{})

		topupBody["amount"] = TopupData.Amount
		topupBody["externalUniqueId"] = TopupData.ID
		topupBody["type"] = TopupData.Type

		// var landingurl = "http://localhost:5000" + "/dashboard/wallet"
		// topupBody["landingUrl"] = "http://localhost:5000/dashboard/wallet"
		if TopupData.Type == "ZA_PEACH_CARD" {
			var TopupCardData TopupCard
			if err := c.ShouldBindJSON(&TopupCardData); err != nil {
				c.JSON(400, gin.H{
					"error":   "Invalid request payload",
					"success": false,
				})
				return
			}

			topupBody["topupCardData"] = map[string]interface{}{
				"accountType":    TopupCardData.AccountType,
				"alias":          TopupCardData.Alias,
				"cardholderName": TopupCardData.CardHolderName,
				"cvv":            TopupCardData.Cvv,
				"dob":            TopupCardData.Dob,
				"expiry":         TopupCardData.Expiry,
				"pan":            TopupCardData.Pan,
			}

		}

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		response, err := UkhesheClient.Post("/wallets/"+utils.ConvertIntToString(int(wallet.Ukheshe_Id))+"/topups", topupBody)

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if response.Status != 200 {
			var responseBody []map[string]interface{}

			if err := json.Unmarshal(response.Data, &responseBody); err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			c.JSON(500, gin.H{
				"success": false,
				"error":   responseBody[0]["description"],
			})
			return
		}

		var responseBody map[string]interface{}

		if err := json.Unmarshal(response.Data, &responseBody); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if err := saveTopupData(responseBody, wallet.Ukheshe_Id, TopupData.ID); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(201, gin.H{
			"success":    true,
			"topup_data": responseBody,
		})
	}
}

func saveTopupData(topupData map[string]interface{}, walletId uint32, topupId uuid.UUID) error {
	var TopupDataBody map[string]interface{}
	var UkhesheClient = configs.MakeAuthenticatedRequest(true)
	fmt.Println(topupData)
	response, err := UkhesheClient.Get("/wallets/+" + utils.ConvertIntToString(int(walletId)) + "/topups/" + utils.ConvertIntToString(int(topupData["topupId"].(float64))))

	if err != nil {
		return err
	}

	if err := json.Unmarshal(response.Data, &TopupDataBody); err != nil {
		return err
	}

	var TopupData models.Topup

	TopupData.Amount = TopupDataBody["amount"].(float64)
	TopupData.TopupType = TopupDataBody["type"].(string)

	createdAt, err := time.Parse(time.RFC3339, TopupDataBody["created"].(string))

	location, err := time.LoadLocation("GMT")

	if err != nil {
		return err
	}

	createdAt = createdAt.In(location)

	TopupData.CreatedAt = createdAt.Unix()
	TopupData.Currency = TopupDataBody["currency"].(string)
	TopupData.SubType = TopupDataBody["subType"].(string)
	TopupData.GateWay = TopupDataBody["gateway"].(string)
	TopupData.GateWayTransactionId = TopupDataBody["gatewayTransactionId"].(string)
	TopupData.TopUpId = uint32(TopupDataBody["topupId"].(float64))

	expiresAt, err := time.Parse(time.RFC3339, TopupDataBody["expires"].(string))

	expiresAt = expiresAt.In(location)
	TopupData.ExpiresAt = expiresAt.Unix()
	TopupData.Ukheshe_Wallet_Id = TopupDataBody["walletId"].(uint32)

	if err := configs.DB.Model(&models.Topup{}).Where("id = ?", topupId).Updates(&TopupData).Error; err != nil {
		return err
	}

	return nil

}

func GetTransactionHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get transaction_type param
		transaction_type := c.Param("transaction_type")
		allowedValues := []string{"transaction", "topup", "withdraw"}
		found := false
		for _, value := range allowedValues {
			if transaction_type == value {
				found = true
				break
			}
		}

		// If the value is not found, reject the request with a "Not Found" response
		if !found {
			c.Redirect(404, "/*")
			return
		}

		page, limit := c.Query("page"), c.Query("limit")
		if page == "" {
			c.JSON(400, gin.H{
				"error":   "Page is required",
				"success": false,
			})
			return
		}
		if limit == "" {
			c.JSON(400, gin.H{
				"error":   "Limit is required",
				"success": false,
			})
			return
		}
		pageInt := utils.ConvertStringToInt(page)
		limitInt := utils.ConvertStringToInt(limit)

		if pageInt <= 0 {
			c.JSON(400, gin.H{
				"error":   "Invalid page numer",
				"success": false,
			})
			return
		}

		if limitInt <= 0 {
			c.JSON(400, gin.H{
				"error":   "Invalid limit numer",
				"success": false,
			})
			return
		}

		offset := utils.GetOffset(pageInt, limitInt)

		var total int64

		if transaction_type == "transaction" {
			var transaction_history []models.Transaction

			if err := configs.DB.Model(&models.Transaction{}).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			if err := configs.DB.Offset(offset).Limit(limitInt).Find(&transaction_history).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			c.JSON(200, gin.H{
				"transaction_history": transaction_history,
				"metadata": map[string]interface{}{
					"limit": limitInt,
					"page":  pageInt,
					"total": total,
				},
				"success": true,
			})

		} else if transaction_type == "topup" {
			var transaction_history []models.Topup

			if err := configs.DB.Model(&models.Topup{}).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			if err := configs.DB.Offset(offset).Limit(limitInt).Find(&transaction_history).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			c.JSON(200, gin.H{
				"transaction_history": transaction_history,
				"metadata": map[string]interface{}{
					"limit": limitInt,
					"page":  pageInt,
					"total": total,
				},
				"success": true,
			})

		} else if transaction_type == "withdraw" {
			var transaction_history []models.Withdraw

			if err := configs.DB.Model(&models.Withdraw{}).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			if err := configs.DB.Offset(offset).Limit(limitInt).Find(&transaction_history).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			c.JSON(200, gin.H{
				"transaction_history": transaction_history,
				"metadata": map[string]interface{}{
					"limit": limitInt,
					"page":  pageInt,
					"total": total,
				},
				"success": true,
			})

		} else {
			c.JSON(500, gin.H{
				"error":   "Something Went Wrong",
				"success": false,
			})
			return
		}

	}
}
