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
	token         string = "eyJraWQiOiIxIiwiYWxnIjoiUFMyNTYifQ.eyJzdWIiOiI2NDA5MTA1MDkwMDg1Iiwic3JjIjoiQXV0aC1SUCIsImNoIjoieTFib1NMbERZbyIsInJvbGVzIjoiIiwic2VzcyI6IjIxOWZhOGQzLWQwMDItNDQwNC05M2YwLTNhZjI2MmIyOWRiOSIsImlwIjoiNDEuMTg2Ljc4LjkiLCJpc3MiOiJodHRwOi8vZWNsaXBzZS1qYXZhLXNhbmRib3gudWtoZXNoZS5yb2NrcyIsImxvY2FsZSI6ImVuLVVTIiwidWlkIjoxNjUzODYsInBvcyI6W3sibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzAxIn0seyJvIjo0OTAyLCJkIjowLCJwIjoiTEVWRUxfMDIifSx7Im8iOjQ5MDIsImQiOjAsInAiOiJMRVZFTF8wMyJ9LHsibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzA0In0seyJvIjo0OTAyLCJkIjowLCJwIjoiTEVWRUxfMDUifSx7Im8iOjQ5MDIsImQiOjAsInAiOiJMRVZFTF8wNiJ9LHsibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzA3In0seyJvIjo0OTAyLCJkIjowLCJwIjoiTEVWRUxfMDgifSx7Im8iOjQ5MDIsImQiOjAsInAiOiJMRVZFTF8wOSJ9LHsibyI6NDkwMiwiZCI6MCwicCI6IkxFVkVMXzEwIn0seyJvIjo0OTAyLCJkIjowLCJwIjoiVEVOQU5UX1NZU1RFTSJ9XSwiZXhwIjoxNjg3MDE5NTc0LCJpYXQiOjE2ODcwMTg2NzQsInRlbmFudCI6NDkwMn0.meqOQyQh6dkD9JVfi_hARwb3a7ylAd5BICXQPVANq6BRq3BzOzmQusQkL-UAJZHGI4K9AgmPOhrDGe9YGHlKz4EJlKSQSgCJ0KCs5RaH21BVRg7jTHA-pQWop8RAcd-r507Hk7nfHdaUf6fZGiR7M-24Y5GEgHjngopBwcYf5arutTSLW2bjr0amfW66ZndlUuPxxuhmH3uk1VEdj1WJNl4fwKQjoNRnyzJaHi9cKeFcdV04wqByuBsSt0l_xElbRupYlMdSLxa5O6rzizqyK2i3X4mGnah_v3miex4Kzd4I4th8Uce2J4OosG43VgoA3iP3yGsUXNLj9YUqt8nbSw"
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
