package main

import (
	"apz-vas/configs"
	"apz-vas/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	_, err := configs.ConnectDb()
	if err != nil {
		panic(err)
	}

	// use the routes

	routes.InitializeRoutes(router)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Welcome to APZ VAS")
	})
	router.Run(":5000")
}
