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
	token         string = "eyJraWQiOiIxIiwiYWxnIjoiUFMyNTYifQ.eyJzdWIiOiI2NDA5MTA1MDkwMDg1Iiwic3JjIjoiQXV0aC1SUCIsImNoIjoiRDFNK29RSzZqciIsInJvbGVzIjoiIiwic2VzcyI6IjIxOWZhOGQzLWQwMDItNDQwNC05M2YwLTNhZjI2MmIyOWRiOSIsImlwIjoiNDEuMTg2Ljc4LjkiLCJpc3MiOiJodHRwOi8vZWNsaXBzZS1qYXZhLXNhbmRib3gudWtoZXNoZS5yb2NrcyIsImxvY2FsZSI6ImVuLVVTIiwidWlkIjoxNjUzODYsInBvcyI6W3sibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzAxIn0seyJvIjo0OTAyLCJkIjowLCJwIjoiTEVWRUxfMDIifSx7Im8iOjQ5MDIsImQiOjAsInAiOiJMRVZFTF8wMyJ9LHsibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzA0In0seyJvIjo0OTAyLCJkIjowLCJwIjoiTEVWRUxfMDUifSx7Im8iOjQ5MDIsImQiOjAsInAiOiJMRVZFTF8wNiJ9LHsibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzA3In0seyJvIjo0OTAyLCJkIjowLCJwIjoiTEVWRUxfMDgifSx7Im8iOjQ5MDIsImQiOjAsInAiOiJMRVZFTF8wOSJ9LHsibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzEwIn0seyJvIjo0OTAyLCJkIjowLCJwIjoiVEVOQU5UX1NZU1RFTSJ9XSwiZXhwIjoxNjg2NTg1MDQ5LCJpYXQiOjE2ODY1ODQxNDksInRlbmFudCI6NDkwMn0.iW_fY03I82nuMOVerB4OymN2DJvbAVuTjRZjcOnum6KIsVSD3t7nmvj4l9WtaB1f5daDtdpypruPu7pz4iMIeAsqjta8kiE27aGJyA0yO4rmfIFDoLCom5dGz1CINdXkjgPEpQx-giJt0S1kffINR3ByGqh-cNdm3LS7-xI2B75Se_0aRH1vb14myF_3osj13DycuYUtcE2pEu3w1xlncQxZ31N5CPPhUTTH5aLE8_zJtoZlTTZYqlaMIOhAVtMkv_iIrJhl6ii85rFSk7wj2954KmshqFb5I5L5ECPNabfRQ-OxAadLgj4KQMzKk3Fl1sRuMI2yCiP5dlQG3eHyBQ"
	tokenMutex    sync.Mutex
	tokenExpires  time.Time
	tokenDuration = 15 * time.Minute
)

func RefreshTokenPeriodically() {
	// Initially, perform the authentication and get the token

	authenticate()

	// Set the timer to refresh the token periodically
	timer := time.NewTimer(tokenDuration)
	for {
		<-timer.C

		// Refresh the token
		authenticate()

		// Reset the timer for the next refresh
		timer.Reset(tokenDuration)
	}
}

func authenticate() {
	for {

		fmt.Println("authenticating...")
		var ukheshe_link = os.Getenv("UKHESHE_LINK")
		// Send a request to renew the token
		resp, err := axios.Post(ukheshe_link+"/eclipse-conductor/rest/v1/authentication/renew", map[string]interface{}{
			"jwt": token,
		}, nil)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if resp.Status != 200 {
			fmt.Println("Failed to renew the token")
			continue
		}
		// Parse the response and update the stored token
		tokenMutex.Lock()

		var responseBody map[string]interface{}

		json.Unmarshal(resp.Data, &responseBody)

		var tokendata = responseBody["headerValue"].(string)
		var tokensplit = strings.Split(tokendata, "Bearer ")[1]
		token = tokensplit
		expiresStr := responseBody["expires"].(string)
		expires, _ := time.Parse(time.RFC3339, expiresStr)

		// Convert the expiration time to GMT
		location, err := time.LoadLocation("GMT")
		if err != nil {
			continue
		}
		expires = expires.In(location)

		tokenExpires = expires
		tokenMutex.Unlock()
		fmt.Println("Authenticated")
		break
	}
}

func MakeAuthenticatedRequest(add_tenant_id bool) *axios.Instance {
	// Retrieve the token and check its expiration
	tokenMutex.Lock()
	currentToken := token
	expires := tokenExpires
	tokenMutex.Unlock()

	// If the token has expired, perform authentication again
	if time.Now().After(expires) {
		authenticate()
		// Retrieve the updated token
		tokenMutex.Lock()
		currentToken = token
		tokenMutex.Unlock()
	}

	// Create an Axios instance
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
