package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"net/url"
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
			if responseBody[i]["mode"] == "SYSTEM" {
				responseBody = append(responseBody[:i], responseBody[i+1:]...)
			}
			if responseBody[i]["mode"] == "PREPAID_CARD" {
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

func UpdateWallet() gin.HandlerFunc {
	return func(c *gin.Context) {

		var wallet models.Wallet

		var walletId = c.Query("walletId")

		if walletId == "" {
			c.JSON(400, gin.H{
				"error":   "Wallet Id is required",
				"success": false,
			})
			return
		}

		if err := c.ShouldBindJSON(&wallet); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		var walletBody = make(map[string]interface{})
		walletBody["description"] = wallet.Description
		walletBody["name"] = wallet.Name
		walletBody["walletTypeId"] = wallet.WalletTypeID

		var walletData = c.MustGet("wallet_data").(map[string]interface{})

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		response, err := UkhesheClient.Put("/wallets/"+utils.ConvertIntToString(int(walletData["walletId"].(float64))), walletBody)

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
			c.JSON(response.Status, gin.H{
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

		if err := configs.DB.Model(&models.Wallet{}).Where("id = ?", walletId).Updates(&wallet).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"wallet":  responseBody,
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
			if err := json.Unmarshal(response.Data, &responseBody); err != nil {
				c.JSON(500, gin.H{
					"success": false,
					"error":   err.Error(),
				})
				return
			}
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

func GetWithdrawalFees() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet = c.MustGet("wallet_data").(map[string]interface{})

		var amount, typedata = c.Query("amount"), c.Query("type")

		if amount == "" {
			c.JSON(400, gin.H{
				"error":   "Amount is required",
				"success": false,
			})
			return
		}

		if typedata == "" {
			c.JSON(400, gin.H{
				"error":   "Type is required",
				"success": false,
			})
			return
		}

		var requestQuery = url.Values{}
		requestQuery.Add("amount", amount)
		requestQuery.Add("type", typedata)

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		var resp, err = UkhesheClient.Get("/wallets/"+utils.ConvertIntToString(int(wallet["walletId"].(float64)))+"/withdrawals/fees", requestQuery)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if resp.Status != 200 {
			var responseBody []map[string]interface{}

			if err := json.Unmarshal(resp.Data, &responseBody); err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			c.JSON(500, gin.H{
				"error":   responseBody[0]["description"],
				"success": false,
			})
			return
		}

		var responseBody map[string]interface{}
		if err := json.Unmarshal(resp.Data, &responseBody); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"fees":    responseBody["feeAmount"],
		})
	}
}

func WithDrawFromWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		// post https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/{walletId}/withdrawals
		var wallet = c.MustGet("wallet_data").(map[string]interface{})
		var requirestBody = c.MustGet("request_body").(map[string]interface{})
		var withdraw models.Withdraw

		var request_body_bytes, errf = json.Marshal(requirestBody)

		if errf != nil {
			c.JSON(500, gin.H{
				"error":   errf.Error(),
				"success": false,
			})
			return

		}

		if err := json.Unmarshal(request_body_bytes, &withdraw); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return

		}

		if withdraw.Amount > wallet["availableBalance"].(float64) {
			c.JSON(500, gin.H{
				"error":   "Insufficient Balance",
				"success": false,
			})
			return
		}

		var organization = c.MustGet("organization_data").(models.Organization)
		withdraw.Location = c.ClientIP()
		withdraw.WalletId = utils.ConvertStringToUUID(wallet["externalUniqueId"].(string))
		if withdraw.DeliveryToPhone == "" {
			withdraw.DeliveryToPhone = organization.Phone_Number1
		}

		withdraw.OrganizationWalletId = wallet["walletId"].(float64)

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
		withdrawBody["location"] = withdraw.Location
		if withdraw.Type == "ZA_NEDBANK_EFT" || withdraw.Type == "ZA_NEDBANK_EFT_IMMEDIATE" {

			if withdraw.AccountName != "" && withdraw.AccountNumber != "" && withdraw.Bank != "" && withdraw.BankCountry != "" && withdraw.BranchCode != "" {
				withdrawBody["accountName"] = withdraw.AccountName
				withdrawBody["accountNumber"] = withdraw.AccountNumber
				withdrawBody["bank"] = withdraw.Bank
				withdrawBody["bankCountry"] = withdraw.BankCountry
				withdrawBody["branchCode"] = withdraw.BranchCode
				if withdraw.Reference != "" {
					withdrawBody["reference"] = withdraw.Reference
				} else {
					withdrawBody["reference"] = "Withdrawal From " + organization.Company_Name + "'s wallet"
				}
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

		// add callbackUrl to be the current url + "/withdraw-callback"
		// withdrawBody["callbackUrl"] = utils.GetFullUrlWithProtocol(c) + "/withdraw-callback"

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

		if err := saveWithdrawData(responseBody); err != nil {
			configs.DB.Where("id = ?", withdraw.ID).Delete(&withdraw)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"success":       true,
			"withdraw_data": responseBody,
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
		topupBody["callbackUrl"] = utils.GetFullUrlWithProtocol(c) + "/topup-callback"

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
		c.JSON(200, gin.H{
			"success":    true,
			"topup_data": responseBody,
		})
	}
}

func UkhesheWithdrawCallBack() gin.HandlerFunc {
	return func(c *gin.Context) {

		var withdrawData map[string]interface{}
		if err := c.ShouldBindJSON(&withdrawData); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		if err := saveWithdrawData(withdrawData); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Withdrawal Callback Received",
		})

	}
}

func UkhesheTopupCallBack() gin.HandlerFunc {
	return func(c *gin.Context) {

		var topupData map[string]interface{}
		if err := c.ShouldBindJSON(&topupData); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		if err := saveTopupDataCall(topupData); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Topup Callback Received",
		})

	}
}

func saveTopupDataCall(topupData map[string]interface{}) error {
	var ourTopupId = topupData["externalUniqueId"]
	var topup models.Topup

	topup.Status = topupData["status"].(string)
	// topup.Amount = topupData["amount"].(float64)
	// topup.Currency = topupData["currency"].(string)
	// topup.GateWay = topupData["gateway"].(string)
	// topup.Type = topupData["type"].(string)
	// topup.SubType = topupData["subType"].(string)
	// topup.TopUpId = topupData["topupId"].(float64)
	// topup.OrganizationWalletId = topupData["walletId"].(float64)
	// topup.PaymentReference = topupData["paymentReference"].(string)
	topup.PaId = topupData["paId"].(string)
	// topup.GateWayTransactionId = topupData["gatewayTransactionId"].(string)4

	if topupData["errorDescription"] != nil {
		topup.ErrorDescription = topupData["errorDescription"].(string)
	}

	if err := configs.DB.Model(&models.Topup{}).Where("id = ?", ourTopupId).Updates(&topup).Error; err != nil {
		return err
	}

	return nil
}

func saveWithdrawData(withdrawData map[string]interface{}) error {
	var ourWithdrawId = withdrawData["externalUniqueId"]
	var withdraw models.Withdraw

	location, err := time.LoadLocation("GMT")

	if err != nil {
		return err
	}

	withdraw.Status = withdrawData["status"].(string)
	// withdraw.Amount = withdrawData["amount"].(float64)
	// withdraw.Currency = withdrawData["currency"].(string)
	withdraw.GateWay = withdrawData["gateway"].(string)
	withdraw.Type = withdrawData["type"].(string)
	withdraw.SubType = withdrawData["subType"].(string)
	withdraw.Fee = withdrawData["fee"].(float64)
	withdraw.WitdrawalId = withdrawData["withdrawalId"].(float64)
	// withdraw.OrganizationWalletId = withdrawData["walletId"].(float64)
	// withdraw.Reference = withdrawData["reference"].(string)
	withdraw.DeliveryToPhone = withdrawData["deliverToPhone"].(string)
	if withdrawData["expires"] != nil {
		expiresAt, err := time.Parse(time.RFC3339, withdrawData["expires"].(string))
		if err != nil {
			return err
		}
		expiresAt = expiresAt.In(location)
		withdraw.ExpiresAt = expiresAt.Unix()
	}

	if withdrawData["created"] != nil {
		createdAt, err := time.Parse(time.RFC3339, withdrawData["created"].(string))
		if err != nil {
			return err
		}
		createdAt = createdAt.In(location)
		withdraw.CreatedAt = createdAt.Unix()
	}

	// if withdrawData["extraInfo"] != nil {
	// 	var extraInfo = withdrawData["extraInfo"].(map[string]interface{})
	// 	withdraw.AccountName = extraInfo["accountName"].(string)
	// 	withdraw.AccountNumber = extraInfo["accountNumber"].(string)
	// 	withdraw.BranchCode = extraInfo["branchCode"].(string)
	// }
	if withdrawData["errorDescription"] != nil {
		withdraw.ErrorDescription = withdrawData["errorDescription"].(string)
	}

	if err := configs.DB.Model(&models.Withdraw{}).Where("id = ?", ourWithdrawId).Updates(&withdraw).Error; err != nil {
		return err
	}
	return nil

}

func saveTopupData(topupData map[string]interface{}, walletId float64, topupId uuid.UUID) error {
	var TopupDataBody map[string]interface{}
	var UkhesheClient = configs.MakeAuthenticatedRequest(true)
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
	if TopupDataBody["errorDescription"] != nil {
		TopupData.ErrorDescription = TopupDataBody["errorDescription"].(string)
	}
	TopupData.TopUpId = TopupDataBody["topupId"].(float64)
	TopupData.OrganizationWalletId = TopupDataBody["walletId"].(float64)

	if TopupDataBody["expires"] != nil {
		expiresAt, err := time.Parse(time.RFC3339, TopupDataBody["expires"].(string))
		if err != nil {
			return err
		}
		expiresAt = expiresAt.In(location)
		TopupData.ExpiresAt = expiresAt.Unix()
	}

	if TopupDataBody["created"] != nil {
		createdAt, err := time.Parse(time.RFC3339, TopupDataBody["created"].(string))
		if err != nil {
			return err
		}
		createdAt = createdAt.In(location)
		TopupData.CreatedAt = createdAt.Unix()
	}

	if err := configs.DB.Model(&models.Topup{}).Where("id = ?", topupId).Updates(&TopupData).Error; err != nil {
		return err
	}

	return nil

}

type TransactionHistory struct {
	Amount      float64   `json:"amount"`
	Rebate      float64   `json:"rebate"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Service     string    `json:"service"`
}

func GetWalletBalances() gin.HandlerFunc {
	return func(c *gin.Context) {
		var wallet_body = c.MustGet("wallet_data").(map[string]interface{})
		var availableBalance = wallet_body["availableBalance"].(float64)
		var currentBalance = wallet_body["currentBalance"].(float64)

		c.JSON(200, gin.H{
			"success": true,
			"message": "Wallet Balance retrieved successfully",
			"balances": map[string]interface{}{

				"available_balance": availableBalance,
				"current_balance":   currentBalance,
			},
		})
	}
}

func GetTransactionHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get transaction_type param

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

		var transaction_history = []models.Transaction{}
		var transactions = []TransactionHistory{}
		var wallet = c.MustGet("wallet_data").(map[string]interface{})

		if err := configs.DB.Model(&models.Transaction{}).Where("organization_wallet_id = ?", int(wallet["walletId"].(float64))).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Select("organization_wallet_id, service_id, amount, currency, created_at, rebate").Where("organization_wallet_id = ?", int(wallet["walletId"].(float64))).Order("created_at desc").Offset(offset).Limit(limitInt).Find(&transaction_history).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		for _, transaction := range transaction_history {
			// delete transaction["id"]
			var Service models.VASService
			if err := configs.DB.Model(&models.VASService{}).Select("id, name").Where("id = ?", transaction.ServiceId).First(&Service).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			transactions = append(transactions, TransactionHistory{
				Amount:      transaction.Amount,
				Rebate:      transaction.Rebate,
				Currency:    transaction.Currency,
				CreatedAt:   time.Unix(transaction.CreatedAt, 0),
				Description: transaction.Description,
				Service:     Service.Name,
			})

		}

		c.JSON(200, gin.H{
			"success":             true,
			"message":             "Transaction history retrieved succesfully.",
			"transaction_history": transactions,
			"metadata": map[string]interface{}{
				"limit": limitInt,
				"page":  pageInt,
				"total": total,
			},
		})

	}
}

func GetWithdrawHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get transaction_type param

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

		var transaction_history []models.Withdraw

		if err := configs.DB.Model(&models.Withdraw{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		if err := configs.DB.Select("type, sub_type, gate_way, fee, amount, currency, delivery_to_phone, bank, branch_code, status, expires_at, created_at, account_number, account_name").Order("created_at desc").Offset(offset).Limit(limitInt).Find(&transaction_history).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":             true,
			"message":             "Withdraw history retrieved succesfully.",
			"transaction_history": transaction_history,
			"metadata": map[string]interface{}{
				"limit": limitInt,
				"page":  pageInt,
				"total": total,
			},
		})

	}
}

func GetTopupHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get transaction_type param

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

		var transaction_history []models.Topup

		if err := configs.DB.Model(&models.Topup{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Select("amount, currency, type, sub_type, gate_way, status, created_at, expires_at").Order("created_at desc").Offset(offset).Limit(limitInt).Find(&transaction_history).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":             true,
			"message":             "Topup history retrieved succesfully.",
			"transaction_history": transaction_history,
			"metadata": map[string]interface{}{
				"limit": limitInt,
				"page":  pageInt,
				"total": total,
			},
		})

	}
}
