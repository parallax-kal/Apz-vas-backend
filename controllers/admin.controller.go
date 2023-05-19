package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"

	"github.com/gin-gonic/gin"
)

func SignupAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var admin models.Admin
		if err := ctx.ShouldBindJSON(&admin); err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		if admin.Name == "" {
			ctx.JSON(400, gin.H{
				"error":   "Name is required",
				"success": false,
			})
			return
		}
		if len(admin.Name) < 3 {
			ctx.JSON(400, gin.H{
				"error":   "Name must be at least 3 characters",
				"success": false,
			})
			return
		}
		// VALIDATE EMAIL
		emailError := utils.ValidateEmail(admin.Email)
		if emailError != nil {
			ctx.JSON(400, gin.H{
				"error":   emailError.Error(),
				"success": false,
			})
			return
		}
		// VALIDATE PASSWORD
		passwordError := utils.ValidatePassword(admin.Password)
		if passwordError != nil {
			ctx.JSON(400, gin.H{
				"error":   passwordError.Error(),
				"success": false,
			})
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(admin.Password)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		admin.Password = hashedPassword
		// Create the user
		// tx := configs.DB.Begin()
		if err := configs.DB.Create(&admin).Error; err != nil {
			// tx.Rollback()
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		// tx.Commit()
		token, err := utils.GenerateToken(
			utils.Data{
				ID: admin.ID,
			},
		)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "Admin created successfully",
			"success": true,
			"token":   token,
		})

	}
}

func LoginAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var admin models.Admin
		if err := c.ShouldBindJSON(&admin); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		// VALIDATE EMAIL
		emailError := utils.ValidateEmail(admin.Email)
		if emailError != nil {
			c.JSON(400, gin.H{
				"error":   emailError.Error(),
				"success": false,
			})
			return
		}
		// VALIDATE PASSWORD
		passwordError := utils.ValidatePassword(admin.Password)
		if passwordError != nil {
			c.JSON(400, gin.H{
				"error":   passwordError.Error(),
				"success": false,
			})
			return
		}
		var givenPassword = admin.Password
		if err := configs.DB.Where("email = ?", admin.Email).First(&admin).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		// Compare the given password with the hashed password that we stored in our database
		if err := utils.ComparePassword(admin.Password, givenPassword); err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(
			utils.Data{
				ID: admin.ID,
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
			"message": "Admin logged in successfully",
			"success": true,
			"token":   token,
		})

	}
}
