package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"github.com/gin-gonic/gin"
)

// CreateOrganization creates a new organization

func SignupOrganization() gin.HandlerFunc {
	// return the handler function and make it async and put there goroutine
	return func(c *gin.Context) {
		// Get the body of the request
		var organization models.Organization
		// Bind the body to the organization var
		if err := c.ShouldBindJSON(&organization); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if organization.Name == "" {
			c.JSON(400, gin.H{
				"error":   "Name is required",
				"success": false,
			})
			return
		}
		if len(organization.Name) < 3 {
			c.JSON(400, gin.H{
				"error":   "Name must be at least 3 characters",
				"success": false,
			})
			return
		}
		// VALIDATE EMAIL
		emailError := utils.ValidateEmail(organization.Email)
		if emailError != nil {
			c.JSON(400, gin.H{
				"error":   emailError.Error(),
				"success": false,
			})
			return
		}
		// VALIDATE PASSWORD

		passwordError := utils.ValidatePassword(organization.Password)
		if passwordError != "" {
			c.JSON(400, gin.H{
				"error":   passwordError,
				"success": false,
			})
			return
		}

		// Hash the password
		hashedPassword, error := utils.HashPassword(organization.Password)
		if error != nil {
			c.JSON(400, gin.H{
				"error":   error.Error(),
				"success": false,
			})
			return
		}

		// Set the hashed password to the organization
		organization.Password = hashedPassword

		// Create the organization
		tx := configs.DB.Begin()
		if err := tx.Select("Name", "Email", "Password", "Status").Create(&organization).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(
			&utils.Data{
				ID:    organization.ID,
				Email: organization.Email,
			},
		)

		if err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		tx.Commit()
		c.JSON(200, gin.H{
			"message": "Organization created successfully",
			"token":   token,
			"success": true,
		})
	}

}

func LoginOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var organization models.Organization
		if err := c.ShouldBindJSON(&organization); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		// VALIDATE EMAIL
		emailError := utils.ValidateEmail(organization.Email)
		if emailError != nil {
			c.JSON(400, gin.H{
				"error":   emailError.Error(),
				"success": false,
			})
			return
		}
		// VALIDATE PASSWORD

		passwordError := utils.ValidatePassword(organization.Password)
		if passwordError != "" {
			c.JSON(400, gin.H{
				"error":   passwordError,
				"success": false,
			})
			return
		}
		// Get the organization from the database
		var givenPassword = organization.Password
		if  err := configs.DB.Where("email = ?", organization.Email).First(&organization).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid credentials",
				"success": false,
			})
			return
		}


		// Check if the password is correct
		if err := utils.CheckPasswordHash(givenPassword, organization.Password); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid credentials",
				"success": false,
			})
			return
		}

		// Generate the JWT token
		token, err := utils.GenerateToken(
			&utils.Data{
				ID:    organization.ID,
				Email: organization.Email,
			},
		)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Logged in successfully",
			"success": true,
			"token":   token,
		})
	}
}
