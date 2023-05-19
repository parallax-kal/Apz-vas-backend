package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

func SignupOrganization() gin.HandlerFunc {
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
		if passwordError != nil {
			c.JSON(400, gin.H{
				"error":   passwordError.Error(),
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
		// tx := configs.DB.Begin()
		if err := configs.DB.Select("Name", "Email", "Password", "Status").Create(&organization).Error; err != nil {
			// tx.Rollback()
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		// tx.Commit()

		token, err := utils.GenerateToken(
			utils.Data{
				ID: organization.ID,
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
		if passwordError != nil {
			c.JSON(400, gin.H{
				"error":   passwordError.Error(),
				"success": false,
			})
			return
		}
		// Get the organization from the database
		var givenPassword = organization.Password
		if err := configs.DB.Where("email = ?", organization.Email).First(&organization).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		// Check if the password is correct
		if err := utils.ComparePassword(givenPassword, organization.Password); err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		// Generate the JWT token
		token, err := utils.GenerateToken(
			utils.Data{
				ID: organization.ID,
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

func CreateOrganization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var organization models.Organization
		if err := ctx.ShouldBindJSON(&organization); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		if organization.Name == "" {
			ctx.JSON(400, gin.H{
				"error":   "Name is required",
				"success": false,
			})
			return
		}
		if len(organization.Name) < 3 {
			ctx.JSON(400, gin.H{
				"error":   "Name must be at least 3 characters",
				"success": false,
			})
			return
		}
		// VALIDATE EMAIL
		emailError := utils.ValidateEmail(organization.Email)
		if emailError != nil {
			ctx.JSON(400, gin.H{
				"error":   emailError.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Create(&organization).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Organization created successfully",
			"success": true,
		})

	}
}

func GetOrganizations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var organizations []models.Organization
		if err := configs.DB.Find(&organizations).Error; err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
		}
		ctx.JSON(200, gin.H{
			"message": "Organizations fetched successfully",
			"success": true,
			"data":    organizations,
		})
	}
}
