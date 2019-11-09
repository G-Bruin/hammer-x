package utils

import (
	"strings"
)

func FindDomain(url string) (domain string) {
	var host string
	if strings.Contains(url, "youtube") {
		host = "youtube"
	}
	return host
}
