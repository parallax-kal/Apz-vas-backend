package middlewares

import (
	"apz-vas/configs"
	"apz-vas/models"
	"apz-vas/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vicanso/go-axios"
	"net/http"
	"net/url"
)

func PasswordChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
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

			var user models.User
			var responseBody map[string]interface{}

			json.Unmarshal(response.Data, &responseBody)

			if err := configs.DB.Where("email = ?", responseBody["email"]).First(&user).Error; err != nil {
				c.JSON(400, gin.H{
					"error":   "You don't have an account with us. Please register",
					"success": false,
				})
				c.Abort()
				return
			}

			pass := utils.ConvertGoogleIdToPassword(responseBody["id"].(string))
			err = utils.ComparePassword(pass, user.Password)

			if err != nil {
				c.JSON(400, gin.H{
					"error":   "You didn't register with google",
					"success": false,
				})
				c.Abort()
				return
			}

			delete(requestData, "google_token")
			c.Set("request_body", requestData)
			c.Next()

		} else {

			var userData = c.MustGet("user_data").(models.User)

			if requestData["password"] == nil {
				c.JSON(400, gin.H{
					"error":   "Enter Password!",
					"success": false,
				})
				c.Abort()
				return
			}

			var err = utils.ComparePassword(requestData["password"].(string), userData.Password)

			if err != nil {
				c.JSON(400, gin.H{
					"error":   "Invalid Password",
					"success": false,
				})
				c.Abort()
				return
			}

			delete(requestData, "password")

			c.Set("request_body", requestData)

			c.Next()
		}
	}
}
