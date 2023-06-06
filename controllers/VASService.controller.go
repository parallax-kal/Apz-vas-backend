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

func GetVASServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var vasServices []models.VASService
		// don't include UpdatedAt, CreatedAt columns
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

		if err := configs.DB.Model(&models.VASService{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Select("id, name, description, rebate, status").Offset(offset).Limit(limitInt).Find(&vasServices).Error; err != nil {
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



type OrganizationSubscribedServices struct {
	models.SubScribedServices
	subscribed bool
}

func OrganizationGetSubscribedServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		organization := c.MustGet("user_data").(models.User)

		// return array of all subscribed services of organization together with unsubscribed ones
		// just make another field which is subscribed, if it is subscribed, then it is true, else false

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

		limitInt := utils.ConvertStringToInt(limit)
		pageInt := utils.ConvertStringToInt(page)

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

		var subScribedServices []models.SubScribedServices
		var vasServices []models.VASService
		var organizationSubscribedServices []OrganizationSubscribedServices

		var total int64

		if err := configs.DB.Model(&models.VASService{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Model(&models.VASService{}).Offset(utils.GetOffset(pageInt, limitInt)).Limit(limitInt).Find(&vasServices).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   "VAS Services not found",
				"success": false,
			})
			return
		}

		if err := configs.DB.Model(&models.SubScribedServices{}).Where("api_key = ?", organization.APIKey).Find(&subScribedServices).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   "SubScribed Services not found",
				"success": false,
			})
			return
		}

		// loop through vasServices and check if it is in subScribedServices
		for _, vasService := range vasServices {
			var subscribed bool
			for _, subScribedService := range subScribedServices {
				if vasService.ID == subScribedService.ServiceId {
					subscribed = true
				}
			}
			organizationSubscribedServices = append(organizationSubscribedServices, OrganizationSubscribedServices{
				SubScribedServices: models.SubScribedServices{
					ServiceId: vasService.ID,
					APIKey:    organization.ID,
				},
				subscribed: subscribed,
			})
		}

		c.JSON(200, gin.H{
			"success":            true,
			"message":            "SubScribed Services retrieved successfully",
			"subScribedServices": organizationSubscribedServices,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
		})

	}
}

func SubScribeService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		organization := ctx.MustGet("organization").(models.User)
		var subScribedService models.SubScribedServices
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
