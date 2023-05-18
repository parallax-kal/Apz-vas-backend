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

	router.Run("localhost:5000")

}
