package configs

import (
	"github.com/vicanso/go-axios"
	"net/http"
	"os"
)

var BlueLabelCleint = getBlueLabelClient()

func getBlueLabelClient() *axios.Instance {

	var blueLabelApiKey = os.Getenv("BLUE_LABEL_API_KEY")

	var blueLabelURL = "https://api.qa.bltelecoms.net"

	// append api key to the url as header

	var blueLabelHeaders = http.Header{
		"Content-Type":       []string{"application/json"},
		"apikey":             []string{blueLabelApiKey},
		"Accept":             []string{"application/json"},
		"Trade-Vend-Channel": []string{"API"},
	}

	var blueLabelClient = axios.NewInstance(
		&axios.InstanceConfig{
			BaseURL: blueLabelURL,
			Headers: blueLabelHeaders,
		},
	)

	return blueLabelClient
}
