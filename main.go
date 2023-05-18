package main

import (
	"apz-vas/configs"
	// "apz-vas/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	_, err := configs.ConnectDb()
	if err != nil {
		panic(err)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to APZ-VAS API",
		})
	})

	router.Run("127.0.0.1:5000")

}
