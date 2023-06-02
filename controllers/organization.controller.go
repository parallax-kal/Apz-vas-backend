package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
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

		org, validationError := ValidateUser(organization)

		if validationError != nil {
			c.JSON(400, gin.H{
				"error":   validationError.Error(),
				"success": false,
			})
			return
		}

		user, err := CreateUser(*org, false)

		if err != nil {
			// check if it is about email existing

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

		user, validationError := ValidateUser(organization)

		if validationError != nil {
			ctx.JSON(400, gin.H{
				"error":   validationError.Error(),
				"success": false,
			})
			return
		}
		_, err := CreateUser(*user, false)
		if err != nil {
			ctx.JSON(500, gin.H{
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
