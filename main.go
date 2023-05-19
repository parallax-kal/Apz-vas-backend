package main

import (
	"apz-vas/configs"
	"apz-vas/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	_, err := configs.ConnectDb()
	if err != nil {
		panic(err)
	}

	// use the routes

	routes.InitializeRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Welcome to APZ-VAS API")
	})
	router.Run("127.0.0.1:5000")
}
