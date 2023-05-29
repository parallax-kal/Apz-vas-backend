package main

import (
	"apz-vas/configs"
	"apz-vas/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	
	_, err := configs.ConnectDb()
	if err != nil {
		panic(err)
	}
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to APZ VAS")
	})

	routes.InitializeRoutes(router)
	router.Run("127.0.0.1:5000")
}
