package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicanso/go-axios"
)

func PasswordChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := c.Request.Header.Get("Authorization")
		if tok == "" {
			c.JSON(401, gin.H{
				"error":   "Invalid Token",
				"success": false,
			})
			c.Abort()
			return
		}
		userData, error := utils.ExtractDataFromToken(tok)

		if error != nil {
			c.JSON(401, gin.H{
				"error":   error.Error(),
				"success": false,
			})
			c.Abort()
			return
		}

		var user models.User

		if err := configs.DB.Select("status, role, name, id", "email").Where("id = ?", userData.ID).First(&user).Error; err != nil {
			c.JSON(401, gin.H{
				"error":   "Invalid token",
				"success": false,
			})
			c.Abort()
			return
		}

		if user.Status != "Active" {
			c.JSON(401, gin.H{
				"error":   "User is not active. Contact Admin to know reason!",
				"success": false,
			})
			c.Abort()
			return
		}

		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid Request Data",
				"success": false,
			})
			c.Abort()
			return
		}

		var tokenString = requestData["google_token"]

		if tokenString != nil {

			var headers = http.Header{
				"Authorization": []string{"Bearer " + tokenString.(string)},
				"Accept":        []string{"application/json"},
			}

			var query = url.Values{
				"access_token": []string{tokenString.(string)},
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
				c.Abort()
				return
			}

			var responseBody map[string]interface{}

			if err := json.Unmarshal(response.Data, &responseBody); err != nil {
				c.JSON(500, gin.H{
					"error":   err.Error(),
					"message": false,
				})
				c.Abort()
				return
			}

			if responseBody["email"] != user.Email {
				c.JSON(400, gin.H{
					"error":   "You didn't register with google",
					"success": false,
				})
				c.Abort()
				return
			}

			if err := configs.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
				c.JSON(400, gin.H{
					"error":   "You don't have an account with us. Please register",
					"success": false,
				})
				c.Abort()
				return
			}

			pass := utils.ConvertGoogleIdToPassword(responseBody["id"].(string))
			err = utils.ComparePassword(pass, strings.Split(user.Passwords, ",")[1])

			if err != nil {
				c.JSON(400, gin.H{
					"error":   "You didn't register with google",
					"success": false,
				})
				c.Abort()
				return
			}

			delete(requestData, "google_token")
			c.Set("user_data", user)
			c.Set("request_body", requestData)
			c.Next()

		} else {

			if requestData["password"] == nil {
				c.JSON(400, gin.H{
					"error":   "Enter Password!",
					"success": false,
				})
				c.Abort()
				return
			}

			if err := configs.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
				c.JSON(400, gin.H{
					"error":   "You don't have an account with us. Please register",
					"success": false,
				})
				c.Abort()
				return
			}

			var err = utils.ComparePassword(requestData["password"].(string), strings.Split(user.Passwords, ",")[0])

			if err != nil {
				c.JSON(400, gin.H{
					"error":   "Invalid Password",
					"success": false,
				})
				c.Abort()
				return
			}

			delete(requestData, "password")
			c.Set("user_data", user)
			c.Set("request_body", requestData)

			c.Next()
		}
	}
}
