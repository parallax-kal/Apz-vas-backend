package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		emailError := utils.ValidateEmail(user.Email)
		if emailError != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid Email address",
				"success": false,
			})
			return
		}
		// VALIDATE PASSWORD

		passwordError := utils.ValidatePassword(user.Password)
		if passwordError != nil {
			c.JSON(400, gin.H{
				"error":   passwordError.Error(),
				"success": false,
			})
			return
		}
		var existingUser models.User
		if err := configs.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		// check has
		err := utils.ComparePassword(user.Password, existingUser.Password)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "Email or password is incorrect",
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(
			utils.Data{
				ID:   existingUser.ID,
				Role: existingUser.Role,
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
			"message": "User logged in successfully",
			"success": true,
			"token":   token,
		})

	}
}

func CreateUser(user models.User) (*models.User, error) {

	if user.Name == "" {
		return nil, errors.New("Name is required.")
	}
	if len(user.Name) < 3 {
		return nil, errors.New("Name must be at least 3 characters")
	}
	// VALIDATE EMAIL
	emailError := utils.ValidateEmail(user.Email)
	if emailError != nil {
		return nil, errors.New(emailError.Error())
	}
	// VALIDATE PASSWORD
	passwordError := utils.ValidatePassword(user.Password)
	if passwordError != nil {

		return nil, passwordError
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {

		return nil, err
	}
	// CREATE USER

	user.Password = hashedPassword

	if err := configs.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}
