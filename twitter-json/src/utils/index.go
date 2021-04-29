package utils

import (
	"regexp"
)

// ReplaceDomain replace twitter domain to cf workers domain
func ReplaceDomain(url string, domain string) string {
	domainReg := regexp.MustCompile(`^((http://)|(https://))?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}(/)`)
	url = domainReg.ReplaceAllString(url, domain+"/")
	return url
}
