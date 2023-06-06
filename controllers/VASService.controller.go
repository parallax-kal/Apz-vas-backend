package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func GetVASServices() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user = c.MustGet("user_data").(models.User)

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
		var vasServicesProviders []VasServiceProvider
		var vasServices []models.VASService

		if user.Role == "Admin" || user.Role == "SuperAdmin" {

			if err := configs.DB.Model(&models.VASService{}).Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			if err := configs.DB.Select("id, name, description, provider_id, rebate, status").Offset(offset).Limit(limitInt).Find(&vasServices).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			for _, vasService := range vasServices {

				var provider models.VASProvider
				if err := configs.DB.Where("id = ?", vasService.ProviderId).First(&provider).Error; err != nil {
					c.JSON(500, gin.H{
						"error":   err.Error(),
						"success": false,
					})
					return
				}
				vasServicesProviders = append(vasServicesProviders, VasServiceProvider{
					vasService,
					provider.Name,
				})
			}
			c.JSON(200, gin.H{
				"message":      "VAS Services retrieved successfully",
				"success":      true,
				"vas_services": vasServicesProviders,
				"metadata": map[string]interface{}{
					"total": total,
					"page":  pageInt,
					"limit": limitInt,
				},
			})
			return

		} else {
			if err := configs.DB.Model(&models.VASService{}).Where("status = ?", "Active").Count(&total).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			// use ProviderService model to get the rebate on a given service and add it to the service struct
			if err := configs.DB.Select("id, name, description, nick_name, rebate, status").Where("status = ?", "Active").Offset(offset).Limit(limitInt).Find(&vasServices).Error; err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}

			c.JSON(200, gin.H{
				"message":      "VAS Services retrieved successfully",
				"success":      true,
				"vas_services": vasServices,
				"metadata": map[string]interface{}{
					"total": total,
					"page":  pageInt,
					"limit": limitInt,
				},
			})
			return
		}
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

func SubScribeService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		organization := ctx.MustGet("organization").(models.User)
		var subScribedService models.SubscribedServices
		if err := ctx.ShouldBindJSON(&subScribedService); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		// type ServiceId is a struct of model.VASServices, but wanna check if it is there
		if subScribedService.ServiceId == uuid.Nil {
			ctx.JSON(400, gin.H{
				"error":   "ServiceId is required",
				"success": false,
			})
			return
		}

		// put organizationId in subScribedService
		subScribedService.APIKey = organization.ID
		if err := configs.DB.Create(&subScribedService).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   "SubScribed Service already exists",
				"success": false,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "SubScribed Service created successfully",
			"success": true,
		})
	}

}
