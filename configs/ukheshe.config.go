package configs

import (
	"encoding/json"
	"fmt"
	"github.com/vicanso/go-axios"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	token         string = os.Getenv("UKHESHE_TOKEN")
	tokenMutex    sync.Mutex
	tokenExpires  time.Time
	tokenDuration = 15 * time.Minute
)

func RefreshTokenPeriodically() {

	authenticate()
	timer := time.NewTimer(tokenDuration)
	for {
		<-timer.C
		authenticate()
		timer.Reset(tokenDuration)
	}
}

func authenticate() {
	for {

		fmt.Println("Authenticating...")
		var ukheshe_link = os.Getenv("UKHESHE_LINK")
		resp, err := axios.Post(ukheshe_link+"/eclipse-conductor/rest/v1/authentication/renew", map[string]interface{}{
			"jwt": token,
		}, nil)

		if err != nil {
			continue
		}

		if resp.Status != 200 {
			fmt.Println("Unable to renew the token.")
			continue
		}
		tokenMutex.Lock()

		var responseBody map[string]interface{}

		json.Unmarshal(resp.Data, &responseBody)

		var tokendata = responseBody["headerValue"].(string)
		var tokensplit = strings.Split(tokendata, "Bearer ")[1]
		token = tokensplit
		expiresStr := responseBody["expires"].(string)
		expires, _ := time.Parse(time.RFC3339, expiresStr)
		tokenExpires = expires
		tokenMutex.Unlock()
		fmt.Println("Authenticated.")
		break
	}
}

func MakeAuthenticatedRequest(add_tenant_id bool) *axios.Instance {
	tokenMutex.Lock()
	currentToken := token
	expires := tokenExpires
	tokenMutex.Unlock()

	if time.Now().After(expires) {
		authenticate()
		tokenMutex.Lock()
		currentToken = token
		tokenMutex.Unlock()
	}

	var ukheshe_link = os.Getenv("UKHESHE_LINK")

	var ukheshe_configs *axios.InstanceConfig

	ukheshe_headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + currentToken},
	}

	var ukheshe_tenat_id = os.Getenv("UKHESHE_TENANT_ID")

	if add_tenant_id {
		ukheshe_configs = &axios.InstanceConfig{
			BaseURL: ukheshe_link + "/eclipse-conductor/rest/v1/tenants/" + ukheshe_tenat_id,
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
