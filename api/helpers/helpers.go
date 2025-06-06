package helpers

import (
	"fmt"
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	if !strings.HasPrefix(url, "http") {
		return fmt.Sprintf("http://%s", url)
	}
	return url
}

// prevent users to use server domain, ex: localhost:9000
func RemoveDomainError(url string) bool {
	appDomain := os.Getenv("APP_DOMAIN")
	// if the url is "exact" localhost:9000
	if url == appDomain {
		return false
	}

	// if the url from user using https, http, wwww.
	newUrl := strings.Replace(url, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	return newUrl != appDomain
}
