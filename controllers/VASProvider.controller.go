	package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

func CreateVASProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasProvider models.VASProvider
		if err := c.ShouldBindJSON(&vasProvider); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
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

		if err := configs.DB.Create(&vasProvider).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "VAS Provider created successfully",
		})
	}
}

func GetVasProviders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var vasProviders []models.VASProvider

		page, limit := ctx.Query("page"), ctx.Query("limit")

		if page == "" {
			ctx.JSON(400, gin.H{
				"error":   "Page is required",
				"success": false,
			})
			return
		}

		if limit == "" {
			ctx.JSON(400, gin.H{
				"error":   "Limit is required",
				"success": false,
			})
			return
		}

		offset := utils.GetOffset(page, limit)

		if err := configs.DB.Offset(offset).Limit(utils.ConvertStringToInt(limit)).Find(&vasProviders).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "VAS Providers fetched successfully",
			"data":    vasProviders,
		})
	}
}

func GetProviderServices() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerService models.ProviderService

		provider_id, page, limit := ctx.Query("provider_id"), ctx.Query("page"), ctx.Query("limit")

		if provider_id == "" {
			ctx.JSON(400, gin.H{
				"error":   "Provider ID is required",
				"success": false,
			})
			return
		}

		if page == "" {
			ctx.JSON(400, gin.H{
				"error":   "Page is required",
				"success": false,
			})
			return
		}

		if limit == "" {
			ctx.JSON(400, gin.H{
				"error":   "Limit is required",
				"success": false,
			})
			return
		}

		offset := utils.GetOffset(page, limit)

		if err := configs.DB.Where("vas_provider_id = ?", provider_id).Offset(offset).Limit(utils.ConvertStringToInt(limit)).Find(&providerService).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"success":           true,
			"message":           "VAS Services of Provider fetched successfully",
			"provider_services": providerService,
		})

	}
}

func GetVasProvidersWithService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var providerServices []models.ProviderService
		page, limit := ctx.Query("page"), ctx.Query("limit")
		if page == "" {
			ctx.JSON(400, gin.H{
				"error":   "Page is required",
				"success": false,
			})
			return
		}

		if limit == "" {
			ctx.JSON(400, gin.H{
				"error":   "Limit is required",
				"success": false,
			})
			return
		}
		offset := utils.GetOffset(page, limit)
		if err := configs.DB.Preload("VASProvider").Preload("VASService").Offset(offset).Limit(utils.ConvertStringToInt(limit)).Find(&providerServices).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"success":           true,
			"message":           "VAS Providers with their Services fetched successfully",
			"provider_services": providerServices,
		})

	}
}

func UpdateVasProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var vasProvider models.VASProvider
		if err := ctx.ShouldBindJSON(&vasProvider); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}
		if vasProvider.Name == "" {
			ctx.JSON(400, gin.H{
				"error":   "Name is required",
				"success": false,
			})
		}
		if len(vasProvider.Name) < 3 {
			ctx.JSON(400, gin.H{
				"error":   "Name must be at least 3 characters",
				"success": false,
			})
		}
		if vasProvider.Description == "" {
			ctx.JSON(400, gin.H{
				"error":   "Description is required",
				"success": false,
			})
		}
		if len(vasProvider.Description) < 3 {
			ctx.JSON(400, gin.H{
				"error":   "Description must be at least 3 characters",
				"success": false,
			})
		}

		if err := configs.DB.Model(&vasProvider).Where("id = ?", vasProvider.ID).Updates(&vasProvider).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "VAS Provider updated successfully",
		})
	}
}

func DeleteVasProvider() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var vasProvider models.VASProvider
		if err := ctx.ShouldBindJSON(&vasProvider); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		if err := configs.DB.Delete(&vasProvider).Where("id = ?", vasProvider.ID).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "VAS Provider deleted successfully",
		})
	}
}

func UpdateProviderService() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var providerService models.ProviderService

		if err := ctx.ShouldBindJSON(&providerService); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		if err := configs.DB.Model(&providerService).Where("ID = ?", providerService.ID).Updates(&providerService).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "VAS Provider Service updated successfully",
		})
	}
}
