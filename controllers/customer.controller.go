package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

func GetCustomers() gin.HandlerFunc {
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
		// get offset

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
		var customers []models.Customer
		organization := c.MustGet("organization_data").(*models.Organization)
		// get the metadata(total)
		var total int64
		if err := configs.DB.Model(&models.Customer{}).Where("organization_id = ?", organization.ID).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Where("api_key = ?", organization.APIKey).Offset(offset).Limit(limitInt).Find(&customers).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success":   true,
			"message":   "Customers retried successfully",
			"customers": customers,
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
		})

	}
}

func CreateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request payload",
				"success": false,
			})
			return
		}

		organization := c.MustGet("user_data").(*models.User)
		customer.OrganizationId = organization.ID
		if err := configs.DB.Create(&customer).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		c.JSON(201, gin.H{
			"message":  "Customer created successfully",
			"customer": customer,
			"success":  true,
		})
	}
}
