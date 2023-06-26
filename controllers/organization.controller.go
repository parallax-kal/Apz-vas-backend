package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func GetYourOrganizationData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organization = c.MustGet("organization_data").(models.Organization)
		c.JSON(200, gin.H{
			"success":      true,
			"message":      "Organization data fetched successfully",
			"organization": organization,
		})
	}
}

func SignupOrganizationContinue() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = c.MustGet("user_data").(models.User)
		var organization models.Organization

		if err := c.ShouldBindJSON(&organization); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		organization.UserId = user.ID
		organization.Email = user.Email
		organization.Owner_Name = user.Name

		if err := configs.DB.Create(&organization).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var organizationBody = make(map[string]interface{})

		organizationBody["email"] = organization.Email
		organizationBody["name"] = organization.Company_Name
		organizationBody["phone1"] = organization.Phone_Number1
		if organization.Phone_Number2 != "" {
			organizationBody["phone2"] = organization.Phone_Number2
		}
		if organization.Tax_Number != "" {

			organizationBody["taxNumber"] = organization.Tax_Number
		}
		if organization.Trading_Name != "" {

			organizationBody["tradingName"] = organization.Trading_Name
		}
		if organization.Company_Number != "" {
			organizationBody["companyNumber"] = organization.Company_Number
		}

		organizationBody["externalUniqueId"] = organization.ID

		if organization.Organization_Type != "" {
			organizationBody["type"] = organization.Organization_Type
		}

		if organization.Industrial_Classification != "" {
			organizationBody["industrialClassification"] = organization.Industrial_Classification
		}

		if organization.Industrial_Sector != "" {
			organizationBody["industrialSector"] = organization.Industrial_Sector
		}

		if organization.Registration_Date != "" {
			organizationBody["businessRegistrationDate"] = organization.Registration_Date // format("20230610")

		}
		if organization.BusinessType != "" {
			organizationBody["businessType"] = organization.BusinessType
		}

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		var response, ukesheResponseError = UkhesheClient.Post("/organisations", organizationBody)

		if ukesheResponseError != nil {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			c.JSON(500, gin.H{
				"success": false,
				"error":   ukesheResponseError.Error(),
			})
			return
		}

		if response.Status != 200 {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			var responseBody []map[string]interface{}
			if err := json.Unmarshal(response.Data, &responseBody); err != nil {
				configs.DB.Where("id = ?", organization.ID).Delete(&organization)
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
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if err := configs.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("status", "Active").Error; err != nil {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var version = responseBody["version"].(float64)
		var organizationId = responseBody["organisationId"].(float64)

		if err := configs.DB.Model(&models.Organization{}).Where("id = ?", organization.ID).Updates(map[string]interface{}{
			"version":    version,
			"ukheshe_id": organizationId,
		}).Error; err != nil {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"data":    responseBody,
		})

	}
}

func SignupOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var usermodel utils.UserEmailedData
		if err := c.ShouldBindJSON(&usermodel); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		validationError := ValidateUser(usermodel, true)

		if validationError != nil {
			c.JSON(400, gin.H{
				"error":   validationError.Error(),
				"success": false,
			})
			return
		}

		usermodel.Role = "Organization"

		token, err := utils.GenerateTokenFromUserData(usermodel)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var link = os.Getenv("FRONTEND_URL") + "/signup/continue?token=" + token

		if err := utils.SendMail(usermodel.Email, "Organization Signup", "Please click on the link below to continue with your registration: "+link); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "We have sent you an email with a link to continue with your registration",
			"success": true,
		})
	}
}

func UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organization = c.MustGet("organization_data").(models.Organization)
		var request_body = c.MustGet("request_body").(map[string]interface{})

		requestBodyBytes, errr := json.Marshal(request_body)

		if errr != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   errr.Error(),
			})
			return
		}

		var org models.Organization

		if errf := json.Unmarshal(requestBodyBytes, &org); errf != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   errf.Error(),
			})
			return
		}
		var ukheshe_client = configs.MakeAuthenticatedRequest(true)

		var organizationBody = make(map[string]interface{})
		if org.Company_Name != "" {
			organizationBody["name"] = org.Company_Name
		} else {
			organizationBody["name"] = nil
		}
		if org.Phone_Number1 != "" {
			organizationBody["phone1"] = org.Phone_Number1
		} else {
			organizationBody["phone1"] = nil
		}
		if org.Phone_Number2 != "" {
			organizationBody["phone2"] = org.Phone_Number2
		} else {
			organizationBody["phone2"] = nil
		}
		if org.Tax_Number != "" {
			organizationBody["taxNumber"] = org.Tax_Number
		} else {
			organizationBody["taxNumber"] = nil
		}
		if org.Trading_Name != "" {
			organizationBody["tradingName"] = org.Trading_Name
		} else {
			organizationBody["tradingName"] = nil
		}
		if org.Company_Number != "" {
			organizationBody["companyNumber"] = org.Company_Number
		} else {
			organizationBody["companyNumber"] = nil
		}

		organizationBody["externalUniqueId"] = organization.ID

		if org.Organization_Type != "" {
			organizationBody["type"] = org.Organization_Type
		} else {
			organizationBody["type"] = nil
		}

		if org.Industrial_Classification != "" {
			organizationBody["industrialClassification"] = org.Industrial_Classification
		} else {
			organizationBody["industrialClassification"] = nil
		}

		if org.Industrial_Sector != "" {
			organizationBody["industrialSector"] = org.Industrial_Sector
		} else {
			organizationBody["industrialSector"] = nil
		}

		if org.Registration_Date != "" {
			organizationBody["businessRegistrationDate"] = org.Registration_Date // format("20230610")

		} else {
			organizationBody["businessRegistrationDate"] = nil
		}
		if org.BusinessType != "" {
			organizationBody["businessType"] = org.BusinessType
		} else {
			organizationBody["businessType"] = nil
		}
		organizationBody["version"] = c.MustGet("organization_version").(float64)
		fmt.Println(organizationBody)

		var response, ukesheResponseError = ukheshe_client.Put("/organisations/"+utils.ConvertIntToString(int(organization.Ukheshe_Id)), organizationBody)

		if ukesheResponseError != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   ukesheResponseError.Error(),
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

		if err := configs.DB.Model(&models.Organization{}).Where("id = ?", organization.ID).Update("ukheshe_id", responseBody["organisationId"]).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if err := configs.DB.Model(&models.Organization{}).Where("id = ?", organization.ID).Updates(org).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":      true,
			"organization": organization,
			"message":      "Organization updated successfully",
		})
	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]interface{}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var ticket = data["ticket"].(string)
		var userData, err = utils.ExtractDataFromUserEmailedDataToken(ticket)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		// check if the email exist

		var user models.User

		if err := configs.DB.Where("email = ?", userData.Email).First(&user).Error; err == nil {
			c.JSON(400, gin.H{
				"error":   "email already taken",
				"success": false,
			})
			return
		}

		if err != nil {
			c.JSON(400, gin.H{
				"error":   "email already taken",
				"success": false,
			})
			return
		}

		dataBytes, err := json.Marshal(data)

		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		user.Passwords = ","
		user.Email = userData.Email
		user.Name = userData.Name
		user.Role = "Organization"

		newUser, err := CreateUser(user)

		if err != nil {
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var organization models.Organization
		organization.Email = newUser.Email
		organization.Owner_Name = newUser.Name

		if err := json.Unmarshal(dataBytes, &organization); err != nil {
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		if err := configs.DB.Create(&organization).Error; err != nil {
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		token, err := utils.GenerateEmailToken(newUser.Email, user.ID)

		if err != nil {
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var link = os.Getenv("FRONTEND_URL") + "/organizations/set-password?token=" + token
		fmt.Println(newUser)
		if err := utils.SendMail(newUser.Email, "Organization Signup", "You've been invited to APZ-VAS click here to continue signing up. link: "+link); err != nil {
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var organizationBody = make(map[string]interface{})

		organizationBody["email"] = organization.Email
		organizationBody["name"] = organization.Company_Name
		organizationBody["phone1"] = organization.Phone_Number1
		if organization.Phone_Number2 != "" {
			organizationBody["phone2"] = organization.Phone_Number2
		}
		if organization.Tax_Number != "" {

			organizationBody["taxNumber"] = organization.Tax_Number
		}
		if organization.Trading_Name != "" {

			organizationBody["tradingName"] = organization.Trading_Name
		}
		if organization.Company_Number != "" {
			organizationBody["companyNumber"] = organization.Company_Number
		}

		organizationBody["externalUniqueId"] = organization.ID

		if organization.Organization_Type != "" {
			organizationBody["type"] = organization.Organization_Type
		}

		if organization.Industrial_Classification != "" {
			organizationBody["industrialClassification"] = organization.Industrial_Classification
		}

		if organization.Industrial_Sector != "" {
			organizationBody["industrialSector"] = organization.Industrial_Sector
		}

		if organization.Registration_Date != "" {
			organizationBody["businessRegistrationDate"] = organization.Registration_Date // format("20230610")

		}
		if organization.BusinessType != "" {
			organizationBody["businessType"] = organization.BusinessType
		}

		var UkhesheClient = configs.MakeAuthenticatedRequest(true)

		var response, ukesheResponseError = UkhesheClient.Post("/organisations", organizationBody)

		if ukesheResponseError != nil {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			c.JSON(500, gin.H{
				"success": false,
				"error":   ukesheResponseError.Error(),
			})
			return
		}

		if response.Status != 200 {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
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
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var version = responseBody["version"].(float64)
		var organizationId = responseBody["organisationId"].(float64)

		if err := configs.DB.Model(&models.Organization{}).Where("id = ?", organization.ID).Updates(map[string]interface{}{
			"version":    version,
			"ukheshe_id": organizationId,
		}).Error; err != nil {
			configs.DB.Where("id = ?", organization.ID).Delete(&organization)
			configs.DB.Where("id = ?", newUser.ID).Delete(&newUser)
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"data":    responseBody,
		})

	}
}

func GetOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organizations []models.Organization
		// get page, limit query param
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
		// get offset
		var total int64

		if err := configs.DB.Model(&models.Organization{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Select("api_key, email, status, created_at, company_name, company_number, trading_name, industrial_sector, industrial_classification, phone_number1, phone_number2, organization_type, business_type, id").Order("created_at DESC").Offset(offset).Limit(limitInt).Find(&organizations).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var orgs []map[string]interface{}

		for _, orga := range organizations {
			// convert it into map
			org := utils.StructToMap(orga)
			org["created_at"] = time.Unix(orga.CreatedAt, 0)

			orgs = append(orgs, org)

		}

		c.JSON(200, gin.H{
			"message":       "Organizations retrieved successfully",
			"organizations": orgs,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
			"success": true,
		})
	}
}
