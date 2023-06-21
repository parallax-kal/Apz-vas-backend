package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

var services = []string{
	"airtime",
	"bundle",
}

func VASServiceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var service_id = c.Query("serviceId")

		if service_id == "" {
			c.JSON(400, gin.H{
				"error":   "Service Id is required",
				"success": false,
			})
			return
		}
		service_uuid := utils.ConvertStringToUUID(service_id)

		var Service models.VASService

		if err := configs.DB.Where("id = ?", service_uuid).First(&Service).Error; err != nil {
			c.JSON(404, gin.H{
				"error":   "Service Not Found",
				"success": false,
			})
			c.Abort()
			return
		}

		c.Set("service_data", Service)

		c.Next()

	}
}

func ServiceProviderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var service = c.MustGet("service_data").(models.VASService)
		var provider models.VASProvider
		if err := configs.DB.Where("id = ?", service.ProviderId).First(&provider).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "This Service doesn't have provider. Contact Admin About this.",
				"success": false,
			})
			c.Abort()
			return
		}
		if provider.Status != "Active" {
			c.JSON(400, gin.H{
				"message": "The Provider for this service is not active.",
				"success": false,
			})
			c.Abort()
		}
		c.Set("service_provider_data", provider)
		c.Next()
	}
}

func NickNameService() gin.HandlerFunc {

	return func(c *gin.Context) {

		var nickname = strings.Split(c.Request.RequestURI, "/")[2]

		if nickname == "" {
			c.Redirect(302, "/not-found")
			return
		}

		// if err := configs.DB

		var Service models.VASService

		if err := configs.DB.Where("nick_name = ?", nickname).First(&Service).Error; err != nil {
			c.JSON(404, gin.H{
				"error":   "Service Not Found. It may have been deleted or It is invalid.",
				"success": false,
			})
			c.Abort()
			return
		}

		if Service.Status != "Active" {
			c.JSON(
				400,
				gin.H{
					"error":   "Service is not active. Contact Admin for a reason!",
					"success": true,
				})
			c.Abort()
			return
		}

		c.Set("service_data", Service)

		c.Next()

	}
}

func CheckSubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organization = c.MustGet("organization_data").(models.Organization)
		var Service = c.MustGet("service_data").(models.VASService)
		var subscribedService models.SubscribedServices

		if err := configs.DB.Where("organization_id = ? AND service_id = ?", organization.ID, Service.ID).First(&subscribedService).Error; err != nil {
			c.JSON(403, gin.H{
				"error":   "You have not subscribed to this service",
				"success": false,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
