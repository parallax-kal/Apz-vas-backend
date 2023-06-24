package controllers

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vicanso/go-axios"
	"net/http"
	"net/url"
	"os"
	"strings"
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

		var userPassword = strings.Split(user.Passwords, ",")[1]

		if userPassword == "" {
			c.JSON(400, gin.H{
				"error":   "Invalid Email or Password!",
				"success": false,
			})
			return
		}

		pass := utils.ConvertGoogleIdToPassword(responseBody["id"].(string))

		// check pass

		err = utils.ComparePassword(pass, userPassword)
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

		user.Passwords = "," + pass

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

		if err := utils.SendMail(user.Email, "Welcome to APZ", "Welcome to APZ. We are glad to have you on board"); err != nil {
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
		var user utils.UserEmailedData
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
				"error":   "Invalid Email or Password!",
				"success": false,
			})
			return
		}
		var userPassword = strings.Split(existingUser.Passwords, ",")[0]

		if userPassword == "" {
			c.JSON(400, gin.H{
				"error":   "Invalid Email or Password!",
				"success": false,
			})
			return
		}
		// check has
		err := utils.ComparePassword(user.Password, userPassword)
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid Email or Password!",
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

func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var emailData map[string]interface{}

		if err := c.ShouldBindJSON(&emailData); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if emailData["email"] == nil {
			c.JSON(400, gin.H{
				"error":   "Email is required.",
				"success": false,
			})
			return
		}

		if err := utils.ValidateEmail(emailData["email"].(string)); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user models.User

		if err := configs.DB.Where("email = ?", emailData["email"].(string)).First(&user).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "User with this email does not exist!",
				"success": false,
			})
			return
		}

		var token, err = utils.GenerateEmailToken(emailData["email"].(string), user.ID)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var link = os.Getenv("FRONTEND_URL") + "/reset-password?token=" + token

		if err := utils.SendMail(emailData["email"].(string), "Reset Password", "Please click the link below to reset your password: \n"+link); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Email sent successfully",
			"success": true,
		})

	}
}

func VerifyBeforeResetingPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token = c.Request.Header.Get("Authorization")
		var emailData, err = utils.ExtractEmailData(token)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user models.User

		if err := configs.DB.Where("email = ?", emailData.Email).First(&user).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid Email!",
				"success": false,
			})
			return
		}

		var ticket, errf = utils.GenerateEmailToken(emailData.Email, user.ID)

		if errf != nil {
			c.JSON(500, gin.H{
				"error":   errf.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Verified successfully",
			"ticket":  ticket,
		})

	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newPasswordData map[string]interface{}

		if err := c.ShouldBindJSON(&newPasswordData); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			c.Abort()
			return
		}

		if newPasswordData["newPassword"] == nil {
			c.JSON(400, gin.H{
				"error":   "New Password is required.",
				"success": false,
			})
			return
		}

		var ticket = c.Request.Header.Get("Authorization")
		var emailData, err = utils.ExtractEmailData(ticket)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := utils.ValidatePassword(newPasswordData["newPassword"].(string)); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			c.Abort()
			return
		}

		var newPass, errf = utils.HashPassword(newPasswordData["newPassword"].(string))

		if errf != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			c.Abort()
			return
		}

		var user models.User

		if err := configs.DB.Where("email = ?", emailData.Email).First(&user).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid Email!",
				"success": false,
			})
			return
		}

		var passwords = strings.Split(user.Passwords, ",")

		passwords[0] = newPass

		var newPasss = strings.Join(passwords, ",")

		if err := configs.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("passwords", newPasss).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid Email!",
				"success": false,
			})
			return
		}

		var newToken, errr = utils.GenerateToken(utils.UserData{
			ID:   user.ID,
			Role: user.Role,
		})

		if errr != nil {
			c.JSON(400, gin.H{
				"error":   errr.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Password reset successfully",
			"token":   newToken,
		})

	}
}

func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var token = c.Request.Header.Get("Authorization")
		var userData, err = utils.ExtractDataFromUserEmailedDataToken(token)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user models.User

		user.Email = userData.Email
		user.Name = userData.Name
		user.Role = userData.Role

		user.Passwords = strings.Join([]string{
			userData.Password,
			"",
		}, ",")

		var userd, errf = CreateUser(user)

		if errf != nil {
			c.JSON(400, gin.H{
				"error":   errf.Error(),
				"success": false,
			})
			return
		}
		var newToken, errr = utils.GenerateToken(utils.UserData{
			ID:   userd.ID,
			Role: userd.Role,
		})

		if errr != nil {
			c.JSON(400, gin.H{
				"error":   errr.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "User verified successfully",
			"token":   newToken,
			"success": true,
		})

	}
}

func ValidateUser(user utils.UserEmailedData, passwordCheck bool) error {
	if user.Name == "" {
		return errors.New("Name is required.")
	}
	if len(user.Name) < 3 {
		return errors.New("Name must be at least 3 characters")
	}

	if len(user.Name) > 25 {
		return errors.New("Name must be at most 25 characters")
	}

	var userExist models.User

	if err := configs.DB.Where("email = ?", user.Email).First(&userExist).Error; err == nil {
		return errors.New("The email has already been taken.")
	}

	// VALIDATE EMAIL
	emailError := utils.ValidateEmail(user.Email)
	if emailError != nil {
		return errors.New(emailError.Error())
	}
	if passwordCheck {
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
	}
	return nil
}

func CreateUser(user models.User) (*models.User, error) {

	var passwords = strings.Split(user.Passwords, ",")

	if passwords[0] != "" {
		newPass, err := utils.HashPassword(passwords[0])
		if err != nil {
			return nil, err
		}
		user.Passwords = newPass + ","

	} else if passwords[1] != "" {

		newPass, err := utils.HashPassword(passwords[1])
		if err != nil {
			return nil, err
		}
		user.Passwords = "," + newPass
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request_body = c.MustGet("request_body").(map[string]interface{})
		if request_body["newPassword"] == nil {
			c.JSON(400, gin.H{
				"error":   "New Password is required.",
				"success": false,
			})
			return
		}
		var newPassword = request_body["newPassword"].(string)
		var newNewpassword, err = utils.HashPassword(newPassword)

		if err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user_data = c.MustGet("user_data").(models.User)

		var passwords = strings.Split(user_data.Passwords, ",")
		passwords[0] = newNewpassword

		var newPass = strings.Join(passwords, ",")

		var user models.User
		user.Passwords = newPass

		if err := configs.DB.Model(&models.User{}).Where("id = ? ", user_data.ID).Updates(&user).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Password updated successfully",
			"success": true,
		})

	}
}

func GetVerificationLink() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request_body = c.MustGet("request_body").(map[string]interface{})

		if request_body["email"] == nil {
			c.JSON(400, gin.H{
				"error":   "Email is required.",
				"success": false,
			})
			return
		}

		var email = request_body["email"].(string)

		var user_data = c.MustGet("user_data").(models.User)

		if email == user_data.Email {
			c.JSON(400, gin.H{
				"error":   "You are already using this email",
				"success": false,
			})
			return
		}

		if err := utils.ValidateEmail(email); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var token, erro = utils.GenerateEmailToken(email, user_data.ID)

		if erro != nil {
			c.JSON(500, gin.H{
				"error":   erro.Error(),
				"success": false,
			})
			return
		}

		link := os.Getenv("FRONTEND_URL") + "/change-email?token=" + token

		if err := utils.SendMail(email, "Change Email", "Please click the link below to change your email: \n"+link); err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Email sent successfully",
			"success": true,
		})

	}
}

func ChangeEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token = c.Request.Header.Get("Authorization")
		var emailData, err = utils.ExtractEmailData(token)

		if err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		var user models.User

		user.Email = emailData.Email

		if err := configs.DB.Model(&models.User{}).Where("id = ?", emailData.ID).Updates(user).Error; err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Email changed successfully",
			"success": false,
		})
	}
}

func AccountSettings() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request_body = c.MustGet("request_body").(map[string]interface{})
		var user_data = c.MustGet("user_data").(models.User)
		if request_body["email"] != user_data.Email {
			c.JSON(400, gin.H{
				"error":   "You can't update email in this way!",
				"success": false,
			})
			c.Abort()
			return
		}
		requestBodyBytes, errr := json.Marshal(request_body)

		if errr != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   errr.Error(),
			})
			return
		}

		var user models.User

		if errf := json.Unmarshal(requestBodyBytes, &user); errf != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   errf.Error(),
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

		if err := configs.DB.Model(&models.User{}).Where("id = ?", user_data.ID).Updates(user).Error; err != nil {
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
	return func(c *gin.Context) {

		page, limit := c.Query("page"), c.Query("limit")

		if page == "" {
			c.JSON(400, gin.H{
				"error":   "Page is required",
				"success": false,
			})
			return
		}
		if limit == "" {
			c.JSON(400, gin.H{
				"error":   "Limit is required",
				"success": false,
			})
			return
		}
		pageInt := utils.ConvertStringToInt(page)
		limitInt := utils.ConvertStringToInt(limit)

		if pageInt <= 0 {
			c.JSON(400, gin.H{
				"error":   "Invalid page numer",
				"success": false,
			})
			return
		}

		if limitInt <= 0 {
			c.JSON(400, gin.H{
				"error":   "Invalid limit numer",
				"success": false,
			})
			return
		}

		offset := utils.GetOffset(pageInt, limitInt)
		// get offset
		var total int64

		var users []models.User

		if err := configs.DB.Model(&models.User{}).Count(&total).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := configs.DB.Order("created_at DESC").Offset(offset).Limit(limitInt).Find(&users).Error; err != nil {
			c.JSON(500, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"users":   users,
			"message": "Users fetched successfully",
			"metadata": map[string]interface{}{
				"total": total,
				"page":  pageInt,
				"limit": limitInt,
			},
		})

	}
}

func ValidateUserOrgData() gin.HandlerFunc {
	return func(c *gin.Context) {

		var userData utils.UserEmailedData

		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(400, gin.H{
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		if err := ValidateUser(userData, false); err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		var ticket, errf = utils.GenerateTokenFromUserData(utils.UserEmailedData{
			Email: userData.Email,
			Role:  "Organization",
			Name:  userData.Name,
		})

		if errf != nil {
			c.JSON(500, gin.H{
				"success": false,
				"error":   errf.Error(),
			})
		}

		c.JSON(200, gin.H{
			"success": true,
			"ticket":  ticket,
			"message": "User Validated",
		})

	}
}
