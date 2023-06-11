package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicanso/go-axios"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user_data").(models.User)
		userMap := utils.StructToMap(user)
		delete(userMap, "password")
		delete(userMap, "updated_at")
		delete(userMap, "created_at")
		delete(userMap, "id")

		c.JSON(200, gin.H{
			"success": true,
			"user":    userMap,
		})
	}
}

func GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		tokenSplit := strings.Split(tokenString, "Bearer ")[1]

		var headers = http.Header{
			"Authorization": []string{"Bearer " + tokenSplit},
			"Accept":        []string{"application/json"},
		}

		var query = url.Values{
			"access_token": []string{tokenSplit},
		}

		instance := axios.NewInstance(&axios.InstanceConfig{
			BaseURL: "https://www.googleapis.com/oauth2/v1/userinfo",
			Headers: headers,
		})

		var response, err = instance.Get("/", query)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user models.User
		var responseBody map[string]interface{}

		json.Unmarshal(response.Data, &responseBody)

		if err := configs.DB.Where("email = ?", responseBody["email"]).First(&user).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "You don't have an account with us. Please register",
				"success": false,
			})
			return
		}

		pass := utils.ConvertGoogleIdToPassword(responseBody["id"].(string))

		// check pass

		err = utils.ComparePassword(pass, user.Password)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "You didn't register with google",
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(
			utils.UserData{
				ID:   user.ID,
				Role: user.Role,
			},
		)

		if err != nil {
			c.JSON(500, gin.H{
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

func GoogleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		tokenSplit := strings.Split(tokenString, "Bearer ")[1]

		var headers = http.Header{
			"Authorization": []string{"Bearer " + tokenSplit},
			"Accept":        []string{"application/json"},
		}

		var query = url.Values{
			"access_token": []string{tokenSplit},
		}

		instance := axios.NewInstance(&axios.InstanceConfig{
			BaseURL: "https://www.googleapis.com/oauth2/v1/userinfo",
			Headers: headers,
		})

		var response, err = instance.Get("/", query)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user models.User
		var responseBody map[string]interface{}

		json.Unmarshal(response.Data, &responseBody)

		user.Email = responseBody["email"].(string)
		user.Name = responseBody["name"].(string)
		var pass = utils.ConvertGoogleIdToPassword(responseBody["id"].(string))
		user.Password = pass

		newUser, errr := CreateUser(user)

		if errr != nil {
			c.JSON(400, gin.H{
				"error":   errr.Error(),
				"success": false,
			})
			return
		}

		token, err := utils.GenerateToken(
			utils.UserData{
				ID:   newUser.ID,
				Role: newUser.Role,
			},
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(201, gin.H{
			"message": "User registered successfully",
			"success": true,
			"token":   token,
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

func ValidateUser(user models.User) error {
	if user.Name == "" {
		return errors.New("Name is required.")
	}
	if len(user.Name) < 3 {
		return errors.New("Name must be at least 3 characters")
	}

	if len(user.Name) > 25 {
		return errors.New("Name must be at most 25 characters")
	}

	// VALIDATE EMAIL
	emailError := utils.ValidateEmail(user.Email)
	if emailError != nil {
		return errors.New(emailError.Error())
	}
	// VALIDATE PASSWORD
	passwordError := utils.ValidatePassword(user.Password)
	if passwordError != nil {

		return passwordError
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {

		return err
	}
	// CREATE USER

	user.Password = hashedPassword
	return nil
}

func CreateUser(user models.User) (*models.User, error) {

	newPass, err := utils.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	user.Password = newPass

	if err := configs.DB.Create(&user).Error; err != nil {
		return nil, err
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

		if len(user.Name) > 25 {
			c.JSON(400, gin.H{
				"error":   "Name must be at most 25 characters",
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
