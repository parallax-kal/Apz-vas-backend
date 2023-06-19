package utils

import (
	"apz-vas/configs"
	"apz-vas/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func PayForService(c *gin.Context) error {
	var wallet = c.MustGet("wallet_data").(map[string]interface{})
	var organization = c.MustGet("organization_data").(models.Organization)
	var service = c.MustGet("service_data").(models.VASService)
	// https://eclipse-java-sandbox.ukheshe.rocks/eclipse-conductor/rest/v1/tenants/{tenantId}/wallets/transfers
	var transfer models.Transaction
	var requestBody = c.MustGet("request_body").(map[string]interface{})
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(requestBodyBytes, &transfer); err != nil {
		return err
	}

	if transfer.Amount <= 0 {
		return errors.New("Amount must be greater than 0")
	}
	if wallet["currentBalance"].(float64) < float64(transfer.Amount) {
		return errors.New("Insufficient Funds")
	}
	transfer.ServiceData = string(requestBodyBytes)
	transfer.Amount = (service.Rebate * transfer.Amount / 100) + transfer.Amount
	transfer.Location = c.ClientIP()
	transfer.ServiceId = service.ID
	transfer.Rebate = service.Rebate
	transfer.Currency = wallet["currency"].(string)
	transfer.ExternalId = organization.ID
	transfer.ApzvasWalletId = uint32(ConvertStringToInt(os.Getenv("APZ_VAS_WALLET_ID")))
	transfer.OrganizationWalletId = wallet["walletId"].(float64)
	transfer.Description = "Payment for " + service.Name
	transfer.WalletId = ConvertStringToUUID(wallet["externalUniqueId"].(string))
	fmt.Println(transfer.Amount)
	if err := configs.DB.Create(&transfer).Error; err != nil {
		return err
	}

	var transferBody = make(map[string]interface{})

	transferBody["ignoreLimits"] = false
	transferBody["location"] = transfer.Location
	transferBody["replyPolicy"] = "WHEN_COMPLETE"
	transferBody["onlyCheck"] = false
	transferBody["fromWalletId"] = transfer.OrganizationWalletId
	transferBody["toWalletId"] = transfer.ApzvasWalletId
	transferBody["description"] = transfer.Description
	transferBody["externalId"] = transfer.ExternalId
	transferBody["externalUniqueId"] = transfer.ID
	transferBody["amount"] = transfer.Amount

	var UkhesheClient = configs.MakeAuthenticatedRequest(true)

	var response, ukesheResponseError = UkhesheClient.Post("/wallets/transfers", transferBody)

	if ukesheResponseError != nil {
		configs.DB.Where("id = ?", transfer.ID).Delete(&transfer)
		return ukesheResponseError
	}

	if response.Status != 204 {
		configs.DB.Where("id = ?", transfer.ID).Delete(&transfer)

		var transferResponse []map[string]interface{}

		json.Unmarshal(response.Data, &transferResponse)
		return errors.New(transferResponse[0]["description"].(string))
	}
	return nil
}
