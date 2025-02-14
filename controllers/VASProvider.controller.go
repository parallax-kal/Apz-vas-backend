package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

func GetVasProviders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasProviders []models.VASProvider

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

		if err := configs.DB.Model(&models.VASProvider{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Offset(offset).Limit(limitInt).Find(&vasProviders).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":       true,
			"message":       "VAS Providers fetched successfully",
			"vas_providers": vasProviders,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
		})
	}
}

func GetProviderServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasServices []models.VASService

		provider_id, page, limit := c.Query("provider_id"), c.Query("page"), c.Query("limit")

		if provider_id == "" {
			c.JSON(400, gin.H{
				"error":   "Provider ID is required",
				"success": false,
			})
			return
		}

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

		if err := configs.DB.Model(&models.VASService{}).Where("provider_id = ?", provider_id).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Where("provider_id = ?", provider_id).Offset(offset).Limit(limitInt).Find(&vasServices).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":           true,
			"message":           "VAS Services of Provider fetched successfully",
			"provider_services": vasServices,
		})

	}
}

func UpdateVasProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasProvider models.VASProvider
		var vasProviderId = c.Query("providerId")
		if err := c.ShouldBindJSON(&vasProvider); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if vasProvider.Status == "" {
			if vasProvider.Name == "" {
				c.JSON(400, gin.H{
					"error":   "Name is required",
					"success": false,
				})
				return
			}
			if len(vasProvider.Name) < 3 {
				c.JSON(400, gin.H{
					"error":   "Name must be at least 3 characters",
					"success": false,
				})
				return
			}
			if vasProvider.Description == "" {
				c.JSON(400, gin.H{
					"error":   "Description is required",
					"success": false,
				})
				return
			}
			if len(vasProvider.Description) < 3 {
				c.JSON(400, gin.H{
					"error":   "Description must be at least 3 characters",
					"success": false,
				})
				return
			}

		} else {
			if vasProvider.Status != "Active" && vasProvider.Status != "Inactive" {
				c.JSON(400, gin.H{
					"error":   "Invalid status",
					"success": false,
				})
				return
			}
		}

		if err := configs.DB.Model(&vasProvider).Where("id = ?", vasProviderId).Updates(&vasProvider).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "VAS Provider updated successfully",
		})
	}
}

func UpdateProviderService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// update

	}
}

func DeleteVasProvider() gin.HandlerFunc {
	return func(c *gin.Context) {

		var vasProviderId = c.Query("providerId")

		if vasProviderId == "" {
			c.JSON(400, gin.H{
				"error":   "Provider ID is required",
				"success": false,
			})
			return
		}

		var vasProvider models.VASProvider

		if err := configs.DB.Where("id = ?", vasProviderId).Delete(&vasProvider).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "VAS Provider deleted successfully",
		})

	}
}
