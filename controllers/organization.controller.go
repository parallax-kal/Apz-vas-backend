package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		user.Role = "Organization"
		newUser, err := CreateUser(user)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}
		var organization models.User
		if err := configs.DB.Create(&organization).Error; err != nil {
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
		organization := ctx.MustGet("organization").(models.Organization)
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

		if subScribedService.OrganizationId == uuid.Nil {
			ctx.JSON(400, gin.H{
				"error":   "OrganizationId is required",
				"success": false,
			})
			return
		}
		// put organizationId in subScribedService
		subScribedService.OrganizationId = organization.ID
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

func OrganizationAccountSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		type OrganizationUpdate struct {
			models.Organization
			changedPassword bool
			Password        string
			NewPassword     string
		}

		var organization OrganizationUpdate
		if err := c.ShouldBindJSON(&organization); err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		// VALIDATE EMAIL
		emailError := utils.ValidateEmail(organization.Email)
		if emailError != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Invalid Email address",
			})
		}
		// VALIDATE PASSWORD
		passwordError := utils.ValidatePassword(organization.Password)
		if passwordError != nil {
			c.JSON(400, gin.H{
				"error":   passwordError.Error(),
				"success": false,
			})
			return
		}
		var newOrganization models.Organization
		if organization.changedPassword {
			newOrganization = models.Organization{
				Name:     organization.Name,
				Email:    organization.Email,
				Password: organization.Password,
			}
		} else {
			// delete NewPassword
			newOrganization = models.Organization{
				Name:  organization.Name,
				Email: organization.Email,
			}
		}
		if err := configs.DB.Model(&newOrganization).Updates(newOrganization).Error; err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Organization updated successfully",
			"success": true,
		})

	}
}
