package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func GetOrganizationYourData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = c.MustGet("user_data").(models.User)

		var organization models.Organization

		if err := configs.DB.Where("user_id = ?", user.ID).First(&organization).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    organization,
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

		if err := configs.DB.Create(&organization).Error; err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var organizationBody = make(map[string]interface{})

		organizationBody["email"] = user.Email
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
		if organization.Bank_Name != "" {
			organizationBody["bankName"] = organization.Bank_Name
		}
		organizationBody["accountNumber"] = organization.Account_Number
		organizationBody["externalUniqueId"] = organization.ID.String()

		if organization.Bank_Name != "" && organization.Account_Number != "" {
			organizationBody["bankDetails"] = []map[string]interface{}{
				{
					"att": "bankName",
					"val": organization.Bank_Name,
				},
				{
					"att": "accountNumber",
					"val": organization.Account_Number,
				},
			}
		}

		if organization.Organization_Type != "" {
			organizationBody["type"] = organization.Organization_Type
		}

		if organization.Registration_Date != "" {
			organizationBody["businessRegistrationDate"] = organization.Registration_Date // format("20230610")

		}
		if organization.BusinessType != "" {
			organizationBody["businessType"] = organization.BusinessType
		}

		var response, ukesheResponseError = configs.UkhesheClient.Post("/organisations", organizationBody)

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

				response, ukesheResponseError = configs.UkhesheClient.Post("/organisations", organizationBody)
				if response.Status != 201 {
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

				// set user status to active
				if err := configs.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("status", "Active").Error; err != nil {
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

			if err := configs.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("status", "Active").Error; err != nil {
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
}

func SignupOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var usermodel models.User
		if err := c.ShouldBindJSON(&usermodel); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		validationError := ValidateUser(usermodel)

		if validationError != nil {
			c.JSON(400, gin.H{
				"error":   validationError.Error(),
				"success": false,
			})
			return
		}

		user, err := CreateUser(usermodel)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		token, err := utils.GenerateToken(utils.UserData{
			ID:   user.ID,
			Role: user.Role,
		})
		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "User registered successfully",
			"success": true,
			"token":   token,
		})
	}
}

func UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		validationError := ValidateUser(user)

		if validationError != nil {
			c.JSON(400, gin.H{
				"error":   validationError.Error(),
				"success": false,
			})
			return
		}
		_, err := CreateUser(user)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var organization models.Organization

		if err := c.ShouldBindJSON(&organization); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		organization.UserId = user.ID

		c.JSON(201, gin.H{
			"message": "Organization created successfully",
			"success": true,
		})

	}
}

func GetOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organizations []models.User
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

		if err := configs.DB.Model(&models.User{}).Where("role = ?", "Organization").Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Where("role = ?", "Organization").Offset(offset).Limit(limitInt).Find(&organizations).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message":       "Organizations retrieved successfully",
			"organizations": organizations,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
			"success": true,
		})
	}
}
