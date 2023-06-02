package middlewares

import (
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

func OrganizationAPIMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apikey := ctx.Request.Header.Get("apz-vas-api-key")
		if apikey == "" {
			ctx.JSON(401, gin.H{
				"error":   "Api Key missing",
				"success": false,
			})
			ctx.Abort()
			return
		}
		APIKey := utils.ConvertStringToUUID(apikey)
		organization, err := utils.CheckApiKey(APIKey)
		if err != nil {
			ctx.JSON(401, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			ctx.Abort()
			return
		}
		print(organization)
		ctx.Set("user_data", organization)
		ctx.Next()

	}
}
