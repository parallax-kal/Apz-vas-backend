package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"time"
	"github.com/gin-gonic/gin"
)

func GetVasServiceData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service = c.MustGet("service_data").(models.VASService)

		c.JSON(200, gin.H{
			"success":     true,
			"message":     "VAS Service retrieved successfully",
			"vas_service": service,
		})

	}
}

func CreateVasService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasService models.VASService
		if err := c.ShouldBindJSON(&vasService); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		if vasService.Name == "" {
			c.JSON(400, gin.H{
				"error":   "Name is required",
				"success": false,
			})
			return
		}
		if len(vasService.Name) < 3 {
			c.JSON(400, gin.H{
				"error":   "Name must be at least 3 characters",
				"success": false,
			})
			return
		}
		if vasService.Description == "" {
			c.JSON(400, gin.H{
				"error":   "Description is required",
				"success": false,
			})
			return
		}
		if len(vasService.Description) < 3 {
			c.JSON(400, gin.H{
				"error":   "Description must be at least 3 characters",
				"success": false,
			})
			return
		}

		if err := configs.DB.Create(&vasService).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "VAS Service created successfully",
		})
	}
}

type VasServiceProvider struct {
	models.VASService
	Provider string `json:"provider"`
}

type VasServiceSubscribed struct {
	models.VASService
	Subscribed bool `json:"subscribed"`
}

func GetAdminVASServices() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		var vas_services []models.VASService
		var vas_services_providers []VasServiceProvider

		if err := configs.DB.Model(&models.VASService{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Select("id, name, description, provider_id, rebate, status").Offset(offset).Limit(limitInt).Find(&vas_services).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		for _, vasService := range vas_services {

			var service_provider models.VASProvider
			if err := configs.DB.Where("id = ?", vasService.ProviderId).First(&service_provider).Error; err != nil {
				vas_services_providers = append(vas_services_providers, VasServiceProvider{
					vasService,
					service_provider.Name,
				})
			} else {
				vas_services_providers = append(vas_services_providers, VasServiceProvider{
					vasService,
					service_provider.Name,
				})
			}
		}

		c.JSON(200, gin.H{
			"message":      "VAS Services retrieved successfully",
			"success":      true,
			"vas_services": vas_services_providers,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
		})

	}
}

func GetOrganizationVASServices() gin.HandlerFunc {
	return func(c *gin.Context) {

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
		var vas_services []models.VASService
		var VasServiceSubscribers []VasServiceSubscribed

		if err := configs.DB.Model(&models.VASService{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   "An error occurred. Please try again or contact admin",
				"success": false,
			})
			return
		}

		if err := configs.DB.Select("id, name, description, provider_id, rebate, status").Offset(offset).Limit(limitInt).Find(&vas_services).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   "An error occurred. Please try again.",
				"success": false,
			})
			return
		}
		var totalSubscribed int64

		var organization = c.MustGet("organization_data").(models.Organization)
		for _, vasService := range vas_services {

			var subscribed_service models.SubscribedServices
			if err := configs.DB.Where("service_id = ? AND organization_id = ?", vasService.ID, organization.ID).First(&subscribed_service).Error; err != nil {
				VasServiceSubscribers = append(VasServiceSubscribers, VasServiceSubscribed{
					vasService,
					false,
				})
			} else {

				VasServiceSubscribers = append(VasServiceSubscribers, VasServiceSubscribed{
					vasService,
					true,
				})
			}
		}

		if err := configs.DB.Model(&models.SubscribedServices{}).Where("organization_id = ?", organization.ID).Count(&totalSubscribed).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   "An error occurred. Please try again or contact admin",
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":      true,
			"message":      "VAS Services retrieved successfully",
			"vas_services": VasServiceSubscribers,
			"metadata": map[string]interface{}{
				"total":           total,
				"page":            pageInt,
				"limit":           limitInt,
				"totalSubscribed": totalSubscribed,
			},
		})

	}
}

func UpdateVasService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasService models.VASService
		if err := c.ShouldBindJSON(&vasService); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Model(&vasService).Where("id = ?", vasService.ID).Updates(&vasService).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "VAS Service not found",
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "VAS Service updated successfully",
		})
	}
}

func DeleteVasService() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasService models.VASService
		if err := c.ShouldBindJSON(&vasService); err != nil {
			c.JSON(400, gin.H{
				"error":   "VAS Service not found",
				"success": false,
			})
			return
		}

		if err := configs.DB.Where("id = ?", vasService.ID).Delete(&vasService).Error; err != nil {
			c.JSON(404, gin.H{
				"error":   "VAS Service not found",
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "VAS Service deleted successfully",
		})
	}
}

func OperationOnService() gin.HandlerFunc {
	return func(c *gin.Context) {
		organization := c.MustGet("organization_data").(models.Organization)
		var subScribedService models.SubscribedServices

		var operation = c.Param("operation")
		var operations = []string{
			"subscribe",
			"unsubscribe",
		}
		found := false
		for _, value := range operations {
			if operation == value {
				found = true
				break
			}
		}

		// If the value is not found, reject the request with a "Not Found" response
		if !found {
			c.Redirect(404, "/*")
			return
		}

		var service = c.MustGet("service_data").(models.VASService)

		subScribedService.ServiceId = service.ID
		if operation == "subscribe" {

			subScribedService.OrganizationId = organization.ID
			subScribedService.Subscription = organization.ID.String() + "-" + service.ID.String()
			if err := configs.DB.Create(&subScribedService).Error; err != nil {
				c.JSON(400, gin.H{
					"error":   "SubScribed Service already exists",
					"success": false,
				})
				return
			}
		} else {
			if err := configs.DB.Where("service_id = ? AND organization_id = ? ", subScribedService.ServiceId, organization.ID).Delete(&subScribedService).Error; err != nil {
				c.JSON(400, gin.H{
					"error":   "SubScribed Service doesn't exists",
					"success": false,
				})
				return
			}
		}
		c.JSON(200, gin.H{
			"message": "SubScribed Service created successfully",
			"success": true,
		})
	}

}

type VasServiceTransaction struct {
	Amount      float64                `json:"amount"`
	Rebate      float64                `json:"rebate"`
	Currency    string                 `json:"currency"`
	ServiceData map[string]interface{} `json:"service_data"`
	CreatedAt   time.Time              `json:"created_at"`
}

func GetVasServiceTransactionHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = c.MustGet("user_data").(models.User)

		page, limit := c.Query("page"), c.Query("limit")

		var service = c.MustGet("service_data").(models.VASService)
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

		var vas_service_transactions []models.Transaction
		var vas_service_transactions_history []VasServiceTransaction

		// check where serviceId is equal to the serviceId
		if user.Role == "Organization" {

			var organization models.Organization

			if err := configs.DB.Where("user_id = ?", user.ID).First(&organization).Error; err != nil {
				c.JSON(401, gin.H{
					"error":   "Unauthorized",
					"success": false,
				})
				c.Abort()
				return
			}

			if err := configs.DB.Model(&models.Transaction{}).Select("service_id, external_id").Where("service_id = ? AND external_id = ?", service.ID, organization.ID).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			if err := configs.DB.Model(&models.Transaction{}).Select("service_id, external_id, service_data, rebate, amount, currency, created_at").Where("service_id = ? AND external_id = ?", service.ID, organization.ID).Order("created_at desc").Offset(offset).Limit(limitInt).Find(&vas_service_transactions).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			for _, transaction := range vas_service_transactions {
				var service_data map[string]interface{}
				if err := json.Unmarshal([]byte(transaction.ServiceData), &service_data); err != nil {
					c.JSON(500, gin.H{
						"error":   err.Error(),
						"success": false,
					})
					return
				}
				delete(service_data, "amount")
				delete(service_data, "device_id")
				vas_service_transactions_history = append(vas_service_transactions_history, VasServiceTransaction{
					Amount:      transaction.Amount,
					Rebate:      transaction.Rebate,
					Currency:    transaction.Currency,
					ServiceData: service_data,
					CreatedAt:   time.Unix(transaction.CreatedAt, 0),
				})
			}

		} else {
			if err := configs.DB.Model(&models.Transaction{}).Where("service_id = ?", service.ID).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			if err := configs.DB.Where("service_id = ?", service.ID).Offset(offset).Limit(limitInt).Find(&vas_service_transactions).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"message":                  "VAS Service Transaction History retrieved successfully.",
			"success":                  true,
			"vas_service_transactions": vas_service_transactions_history,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
		},
		)
	}
}
