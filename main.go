package main

import (
	"apz-vas/configs"
	// "apz-vas/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// initialize routes
	// routes.InitializeRoutes(router)
	_, err := configs.ConnectDb()
	if err != nil {
		panic(err)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to APZ-VAS API",
		})
	})

	router.Run(":5000")

}
