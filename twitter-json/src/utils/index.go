package utils

import (
	"os"
	"regexp"
)

// ReplaceDomain replace twitter domain to cf workers domain
func ReplaceDomain(url string) string {
	domainReg := regexp.MustCompile(`^((http://)|(https://))?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}(/)`)
	domain := os.Getenv("CF_WORKERS_DOMAIN")
	url = domainReg.ReplaceAllString(url, domain+"/")
	return url
}
