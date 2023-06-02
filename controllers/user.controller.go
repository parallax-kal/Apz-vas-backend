package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"errors"
	"github.com/gin-gonic/gin"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user_data").(models.User)

		userMap := utils.StructToMap(user)
		delete(userMap, "password")
		delete(userMap, "updated_at")
		delete(userMap, "created_at")
		delete(userMap, "status")
		delete(userMap, "id")

		if user.Role == "Admin" {
			delete(userMap, "api_key")
		}
		c.JSON(200, gin.H{
			"success": true,
			"user":    userMap,
		})
	}
}

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
			utils.UserData{
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

func ValidateUser(user models.User) (*models.User, error) {
	if user.Name == "" {
		return nil, errors.New("Name is required.")
	}
	if len(user.Name) < 3 {
		return nil, errors.New("Name must be at least 3 characters")
	}

	if len(user.Name) > 15 {
		return nil, errors.New("Name must be at most 15 characters")
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
	return &user, nil
}

func CreateUser(user models.User, admin bool) (*models.User, error) {

	// if err := configs.DB.Create(&user).Error; err != nil {
	// 	return nil, err
	// }
	// check if he is an admin or not and give him an api key if he is not admin
	if admin == false {
		// user model has a default api key generation allow it now
		if err := configs.DB.Create(&user).Error; err != nil {
			return nil, err
		}

	} else {
		user.Role = "Admin"
		// check if the user is an admin
		if err := configs.DB.Select("name", "email", "password", "role").Create(&user).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil

}

type changePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	ChangePass  bool   `json:"change_pass"`
}

func AccountSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		if user.Name == "" {
			c.JSON(400, gin.H{
				"error":   "Name is required.",
				"success": false,
			})
			return
		}
		if len(user.Name) < 3 {
			c.JSON(400, gin.H{
				"error":   "Name must be at least 3 characters",
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

		var passwordBody changePassword
		if err := c.ShouldBindJSON(&passwordBody); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		if passwordBody.ChangePass {
			// VALIDATE PASSWORD
			passwordError := utils.ValidatePassword(passwordBody.NewPassword)
			if passwordError != nil {
				c.JSON(400, gin.H{
					"error":   passwordError.Error(),
					"success": false,
				})
				return
			}
			// Hash the password
			hashedPassword, err := utils.HashPassword(passwordBody.NewPassword)
			if err != nil {
				c.JSON(400, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			user.Password = hashedPassword
		}

		if err := configs.DB.Model(&user).Updates(user).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "User updated successfully",
			"success": true,
		})

	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
