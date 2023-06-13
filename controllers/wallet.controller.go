package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"
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
		walletBody["status"] = "ACTIVE"
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
		
		var wallet_body = c.MustGet("wallet_data").(map[string]interface{})

		c.JSON(200, gin.H{
			"success": true,
			"wallet":  wallet_body,
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

func WithDrawFromWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// post https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/{walletId}/withdrawals
		var wallet = c.MustGet("wallet_data").(map[string]interface{})
		var withdraw models.Withdraw

		if err := c.ShouldBindJSON(&withdraw); err != nil {
			fmt.Println(err.Error())
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		var organization = c.MustGet("organization_data").(models.Organization)

		withdraw.WalletId = utils.ConvertStringToUUID(wallet["externalUniqueId"].(string))
		if withdraw.DeliveryToPhone == "" {
			withdraw.DeliveryToPhone = organization.Phone_Number1
		}
		if err := configs.DB.Create(&withdraw).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var withdrawBody = make(map[string]interface{})

		withdrawBody["amount"] = withdraw.Amount
		withdrawBody["externalUniqueId"] = withdraw.ID
		withdrawBody["location"] = c.ClientIP()
		if withdraw.Type == "ZA_NEDBANK_EFT" || withdraw.Type == "ZA_NEDBANK_EFT_IMMEDIATE" {
			if withdraw.AccountName != "" && withdraw.AccountNumber != "" && withdraw.Bank != "" && withdraw.BankCountry != "" && withdraw.BranchCode != "" {
				withdrawBody["accountName"] = withdraw.AccountName
				withdrawBody["accountNumber"] = withdraw.AccountNumber
				withdrawBody["bank"] = withdraw.Bank
				withdrawBody["bankCountry"] = withdraw.BankCountry
				withdrawBody["branchCode"] = withdraw.BranchCode
			} else {
				configs.DB.Where("id = ?", withdraw.ID).Delete(&withdraw)
				c.JSON(500, gin.H{
					"error":   "Enter Bank Account Data",
					"success": false,
				})
				return
			}
		}
		withdrawBody["type"] = withdraw.Type
		withdrawBody["description"] = "Withdrawing From " + organization.Company_Name + "'s wallet"
		withdrawBody["deliverToPhone"] = withdraw.DeliveryToPhone

		var Ukheshe_Client = configs.MakeAuthenticatedRequest(true)
		var response, err = Ukheshe_Client.Post("/wallets/"+utils.ConvertIntToString(int(wallet["walletId"].(float64)))+"/withdrawals", withdrawBody)

		if err != nil {
			configs.DB.Where("id = ?", withdraw.ID).Delete(&withdraw)
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if response.Status != 200 {
			var responseBody []map[string]interface{}

			configs.DB.Where("id = ?", withdraw.ID).Delete(&withdraw)

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
			configs.DB.Where("id = ?", withdraw.ID).Delete(&withdraw)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// if err := saveTopupData(responseBody, wallet["walletId"].(float64), withdraw.ID); err != nil {
		// 	configs.DB.Where("id = ?", withdraw.ID).Delete(&withdraw)
		// 	c.JSON(500, gin.H{
		// 		"success": false,
		// 		"error":   err.Error(),
		// 	})
		// 	return
		// }
		c.JSON(201, gin.H{
			"success":    true,
			"topup_data": responseBody,
		})
	}
}

func TopUpWallet() gin.HandlerFunc {

	return func(c *gin.Context) {
		var TopupData models.Topup

		if err := c.ShouldBindJSON(&TopupData); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		wallet := c.MustGet("wallet_data").(map[string]interface{})

		TopupData.WalletId = utils.ConvertStringToUUID(wallet["externalUniqueId"].(string))

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

		response, err := UkhesheClient.Post("/wallets/"+utils.ConvertIntToString(int(wallet["walletId"].(float64)))+"/topups", topupBody)

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

			configs.DB.Where("id = ?", TopupData.ID).Delete(&TopupData)

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
			configs.DB.Where("id = ?", TopupData.ID).Delete(&TopupData)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if err := saveTopupData(responseBody, wallet["walletId"].(float64), TopupData.ID); err != nil {
			configs.DB.Where("id = ?", TopupData.ID).Delete(&TopupData)
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

func saveTopupData(topupData map[string]interface{}, walletId float64, topupId uuid.UUID) error {
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
	if TopupDataBody["expires"] != nil {
		expiresAt, err := time.Parse(time.RFC3339, TopupDataBody["expires"].(string))
		if err != nil {
			return err
		}
		expiresAt = expiresAt.In(location)
		TopupData.ExpiresAt = expiresAt.Unix()
		TopupData.Ukheshe_Wallet_Id = TopupDataBody["walletId"].(float64)
	}

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
			var transaction_history []map[string]interface{}

			// https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/{walletId}/transactions

			wallet := c.MustGet("wallet_data").(map[string]interface{})

			var UkhesheClient = configs.MakeAuthenticatedRequest(true)

			response, err := UkhesheClient.Get("/wallets/" + utils.ConvertIntToString(int(wallet["walletId"].(float64))) + "/transactions?limit=" + limit + "&offset=" + utils.ConvertIntToString(offset))

			if err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			if response.Status != 200 {
				c.JSON(500, gin.H{
					"error":   "Something went wrong",
					"success": false,
				})
				return
			}

			if err := configs.DB.Model(&models.Transaction{}).Count(&total).Error; err != nil {
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
			var transaction_history []map[string]interface{}

			if err := configs.DB.Model(&models.Topup{}).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			// https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/{walletId}/topups

			wallet := c.MustGet("wallet_data").(map[string]interface{})
			var UkhesheClient = configs.MakeAuthenticatedRequest(true)

			response, err := UkhesheClient.Get("/wallets/" + utils.ConvertIntToString(int(wallet["walletId"].(float64))) + "/topups?limit=" + limit + "&offset=" + utils.ConvertIntToString(offset))

			if err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			if response.Status != 200 {
				c.JSON(500, gin.H{
					"error":   "Something went wrong",
					"success": false,
				})
				return
			}

			if err := json.Unmarshal(response.Data, &transaction_history); err != nil {
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
