package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetWalletTypes() gin.HandlerFunc {
	return func(c *gin.Context) {

		response, err := configs.UkhesheClient.Get("/wallet-types")

		if err != nil {
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

			var expired = configs.CheckTokenExpiry(responseBody[0])
			if expired {
				var err = configs.RenewUkhesheToken()

				if err != nil {
					c.JSON(500, gin.H{
						"success": false,
						"error":   err.Error(),
					})
					return
				}

				response, err := configs.UkhesheClient.Get("/wallet-types")
				if err != nil {
					c.JSON(500, gin.H{
						"success": false,
						"error":   err.Error(),
					})
					return
				}
				if response.Status != 200 {
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

				// fmt.Println(responseBody)

				c.JSON(200, gin.H{
					"success":      true,
					"wallet_types": responseBody,
				})

			} else {
				c.JSON(response.Status, gin.H{
					"success": false,
					"error":   responseBody[0]["description"],
				})
				return
			}

		} else {

			var responseBody []map[string]interface{}

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

		response, err := configs.UkhesheClient.Post("/organisations/"+utils.ConvertIntToString(int(organization.Ukheshe_Id))+"/wallets", walletBody)

		if err != nil {
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

			var expired = configs.CheckTokenExpiry(responseBody[0])
			if expired {
				var err = configs.RenewUkhesheToken()

				if err != nil {
					c.JSON(500, gin.H{
						"success": false,
						"error":   err.Error(),
					})
					return
				}

				response, err := configs.UkhesheClient.Post("/organisations/"+utils.ConvertIntToString(int(organization.Ukheshe_Id))+"/wallets", walletBody)
				if err != nil {
					c.JSON(500, gin.H{
						"success": false,
						"error":   err.Error(),
					})
					return
				}
				if response.Status != 200 {
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

				// fmt.Println(responseBody)

				c.JSON(200, gin.H{
					"success": true,
					"wallet":  responseBody,
				})

			} else {
				c.JSON(response.Status, gin.H{
					"success": false,
					"error":   responseBody[0]["description"],
				})
				return
			}

		} else {

			var responseBody map[string]interface{}

			if err := json.Unmarshal(response.Data, &responseBody); err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}

			// update wallet
			if err := configs.DB.Model(&models.Wallet{}).Where("id = ?", wallet.ID).Update("ukheshe_id", responseBody["walletId"]).Error; err != nil {
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

		fmt.Println(TopupData)

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
		response, err := configs.UkhesheClient.Post("/wallets/"+utils.ConvertIntToString(int(wallet.Ukheshe_Id))+"/topups", topupBody)

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

			var expired = configs.CheckTokenExpiry(responseBody[0])
			if expired {
				var err = configs.RenewUkhesheToken()

				if err != nil {
					c.JSON(500, gin.H{
						"success": false,
						"error":   err.Error(),
					})
					return
				}

				response, err := configs.UkhesheClient.Post("/wallets/"+utils.ConvertIntToString(int(wallet.Ukheshe_Id))+"/topups", topupBody)
				if err != nil {
					c.JSON(500, gin.H{
						"success": false,
						"error":   err.Error(),
					})
					return
				}
				if response.Status != 200 {
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

				fmt.Println(responseBody)

				c.JSON(20, gin.H{
					"success":    true,
					"topup_data": responseBody,
				})

			} else {
				c.JSON(response.Status, gin.H{
					"success": false,
					"error":   responseBody[0]["description"],
				})
				return
			}

		} else {
			var responseBody map[string]interface{}

			if err := json.Unmarshal(response.Data, &responseBody); err != nil {
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

}
