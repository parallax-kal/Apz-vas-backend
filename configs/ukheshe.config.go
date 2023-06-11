package configs

import (
	"encoding/json"
	"fmt"
	"github.com/vicanso/go-axios"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var UkhesheClient = getUkhesheClient(true, true)

type UkhesheTokenData struct {
	Token string `json:"token"`
}

func getToken() string {
	jsonFile, err := os.Open("token.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened token.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var tokenData UkhesheTokenData

	json.Unmarshal(byteValue, &tokenData)

	return tokenData.Token
}

func saveToken(token string) {
	var tokenData = UkhesheTokenData{
		Token: token,
	}

	file, _ := json.MarshalIndent(tokenData, "", " ")

	_ = ioutil.WriteFile("token.json", file, 0644)
}

func getUkhesheClient(useToken bool, add_tenant_id bool) *axios.Instance {
	var ukheshe_link = os.Getenv("UKHESHE_LINK")

	var ukheshe_headers http.Header

	if useToken {
		var Ukheshe_TOKEN = getToken()
		ukheshe_headers = http.Header{
			"Content-Type":  []string{"application/json"},
			"Accept":        []string{"application/json"},
			"Authorization": []string{"Bearer " + Ukheshe_TOKEN},
		}
	} else {
		ukheshe_headers = http.Header{
			"Content-Type": []string{"application/json"},
			"Accept":       []string{"application/json"},
		}
	}

	var ukheshe_tenat_id = os.Getenv("UKHESHE_TENANT_ID")

	var ukheshe_configs *axios.InstanceConfig

	if add_tenant_id {
		ukheshe_configs = &axios.InstanceConfig{
			BaseURL: ukheshe_link + "/eclipse-conductor/rest/v1" + "/tenants/" + ukheshe_tenat_id,
			Headers: ukheshe_headers,
		}
	} else {
		ukheshe_configs = &axios.InstanceConfig{
			BaseURL: ukheshe_link + "/eclipse-conductor/rest/v1",
			Headers: ukheshe_headers,
		}
	}

	var ukheshe_client = axios.NewInstance(
		ukheshe_configs,
	)

	return ukheshe_client

}

func CheckTokenExpiry(response map[string]interface{}) bool {
	if response["code"] == "ARCH014" {
		return true
	}
	return false
}

func RenewUkhesheToken() error {

	fmt.Println("Renewing Ukheshe token")
	var UkhesheClient = getUkhesheClient(false, false)

	var Ukheshe_TOKEN = getToken()

	var body = map[string]string{
		"jwt": Ukheshe_TOKEN,
	}

	var response, err = UkhesheClient.Post("/authentication/renew", body)
	if err != nil {
		return err
	}

	var responseBody map[string]interface{}

	json.Unmarshal(response.Data, &responseBody)

	var token = responseBody["headerValue"].(string)
	var tokensplit = strings.Split(token, "Bearer ")[1]

	saveToken(tokensplit)

	fmt.Println("Ukheshe token renewed")
	return nil
}
