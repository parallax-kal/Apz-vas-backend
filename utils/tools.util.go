package utils

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ConvertGoogleIdToPassword(googleId string) string {
	return os.Getenv("GOOGLE_ID_PREFIX") + googleId + os.Getenv("GOOGLE_ID_SUFFIX")
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func ConvertStringToInt(value string) int {
	convertedValue, _ := strconv.Atoi(value)
	return convertedValue
}

func ConvertIntToString(value int) string {
	return strconv.Itoa(value)
}
func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}

func StructToMap(obj interface{}) map[string]interface{} {
	var data map[string]interface{}
	inrec, _ := json.Marshal(obj)
	json.Unmarshal(inrec, &data)
	return data

}


func GetFullUrlWithProtocol(c *gin.Context) string {
	protocol := "http"
	if c.Request.TLS != nil {
		protocol = "https"
	}
	return protocol + "://" + c.Request.Host + c.Request.RequestURI
}