package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"github.com/gin-gonic/gin"
)


func GetWalletInfo() gin.HandlerFunc{
	return func(c*gin.Context) {
		var wallet models.Wallet

		organization := c.MustGet("user_data").(*models.User)

		if err := configs.DB.Where("organization_id = ?", organization.ID).First(&wallet).Error; err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data": wallet,
		})
	}
}

func CreateWallet() gin.HandlerFunc{
	return func(c*gin.Context) {
		var wallet models.Wallet

		organization := c.MustGet("user_data").(*models.User)

		wallet.OrganizationId = organization.ID
		
		if err := configs.DB.Create(&wallet).Error; err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
				"success": false,
			})
			return
		}
	}
}