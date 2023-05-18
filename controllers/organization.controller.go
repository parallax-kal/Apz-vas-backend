package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"github.com/gin-gonic/gin"

)

// CreateOrganization creates a new organization

func CreateOrganization() gin.HandlerFunc {
	// return the handler function and make it async and put there goroutine
	return func(c *gin.Context) {
		// Get the DB connection
		db, err := configs.ConnectDb()
		if err != nil {
			panic(err)
		}
		// Get the JSON data from the request
		var data models.Organization
		c.BindJSON(&data)
		// Create the organization
		go db.Create(&data)
		// Return the organization
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Organization created successfully!",
			"data":    data,
		})
	}
}
