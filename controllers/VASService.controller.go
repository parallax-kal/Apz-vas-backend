package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"github.com/gin-gonic/gin"
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

func GetVASServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasServices []models.VASService
		// don't include UpdatedAt, CreatedAt columns

		if err := configs.DB.Select("id, name, description, status").Find(&vasServices).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		c.JSON(200, gin.H{
			"success":      true,
			"vas_services": vasServices,
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
			c.JSON(400, gin.H{
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
