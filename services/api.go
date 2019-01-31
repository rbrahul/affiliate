package services

import (
	"fmt"
	"net/http"
)

//DomainMapping - List of amazon affiliate store
var DomainMapping map[string]string

func init() {
	DomainMapping = map[string]string{
		"us": "com",
		"de": "de",
		"uk": "co.uk",
		"it": "it",
		"es": "es",
		"fr": "fr",
		"cn": "cn",
		"jp": "co.jp",
		"in": "in",
		"ca": "ca",
		"br": "com.br",
		"mx": "com.mx",
	}
}

func getAmazonURL(asin string, countryCode string) string {
	domain, ok := DomainMapping[countryCode]
	if !ok {
		fmt.Println("Domain Not found")
		return ""
	}
	return fmt.Sprintf("http://amazon.%s/dp/%s/", domain, asin)
}

// IsProductAvailable - http call to amazon api based on country code and asin
func IsProductAvailable(asin string, countryCode string) bool {
	url := getAmazonURL(asin, countryCode)
	if url == "" {
		fmt.Printf("No Domain matched on that country %s", countryCode)
	} else {
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("ERROR-AMAZON-API: ", response.StatusCode)
		}
		fmt.Println("| ", asin, " | ", countryCode, " | ", response.StatusCode, " |")
		defer response.Body.Close()
		return response.StatusCode == 200
	}
	return false
}
