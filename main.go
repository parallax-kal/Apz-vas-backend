package main

import (
	"apz-vas/configs"
	"apz-vas/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	time.Local = time.UTC
	go configs.RefreshTokenPeriodically()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "apz-vas-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// check not found routes
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Route Not found",
			"success": false,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to APZ VAS")
	})

	routes.InitializeRoutes(router)

	// router.SetTrustedProxies(nil)

	router.Run(":5000")
}
