package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SignupOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organization models.User
		if err := c.ShouldBindJSON(&organization); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		user, err := CreateUser(organization, false)
		if err != nil {
			c.JSON(400, gin.H{
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
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Organization created successfully",
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
	return func(ctx *gin.Context) {
		var organization models.User
		if err := ctx.ShouldBindJSON(&organization); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		_, err := CreateUser(organization, false)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Organization created successfully",
			"success": true,
		})

	}
}

func GetOrganizations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var organizations []models.User
		// get page, limit query param
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
		// get offset
		offset := utils.GetOffset(page, limit)

		if err := configs.DB.Where("role = ?", "Organization").Offset(offset).Limit(utils.ConvertStringToInt(limit)).Find(&organizations).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message":       "Organizations retrieved successfully",
			"organizations": organizations,
			"success":       true,
		})
	}
}

func GetOrganizationSubScribedServices() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get organization from context
		organization := ctx.MustGet("organization").(models.User)
		var subScribedServices []models.SubScribedServices
		// get page, limit query param
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
		// get offset
		offset := utils.GetOffset(page, limit)
		// get subScribedServices
		if err := configs.DB.Where("organization_id = ?", organization.ID).Offset(offset).Limit(utils.ConvertStringToInt(limit)).Find(&subScribedServices).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message":            "SubScribed Services retrieved successfully",
			"subScribedServices": subScribedServices,
			"success":            true,
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
