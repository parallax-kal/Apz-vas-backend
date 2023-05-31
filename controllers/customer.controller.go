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
		offset := utils.GetOffset(page, limit)
		var customers []models.Customer
		organization := c.MustGet("user").(models.User)

		if err := configs.DB.Where("api_key = ?", organization.APIKey).Offset(offset).Limit(utils.ConvertStringToInt(limit)).Find(&customers).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}

		c.JSON(200, gin.H{
			"message":   "Customers retried successfully",
			"customers": customers,
			"success":   true,
		})

	}
}
