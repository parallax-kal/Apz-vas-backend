package controllers

import (
	"github.com/gin-gonic/gin"
	"os"
)

func SignupAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var clientOrigin = ctx.GetHeader("Origin")

		if clientOrigin != os.Getenv("FRONTEND_URL") {
			ctx.JSON(400, gin.H{
				"error":   "Invalid Origin",
				"success": false,
			})
			return
		}

	}
}

func GetAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func DeleteAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func CreateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
