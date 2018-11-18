package main

import (
	"net/url"
	"strings"
)

// SplitUrl just return a splitted string.
func SplitUrl(u string) []string {
	return strings.Split(u, "/")
}

// UrlSource will retrieve the url hostname.
func UrlSource(u string) (string, error) {
	parsedUrl, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	return parsedUrl.Hostname(), nil
}

// IsUrlValid will exclude those url containing `.gif` and `logo`.
func IsUrlValid(url string) bool {
	return !strings.Contains(url, ".gif") && !strings.Contains(url, "logo") && !strings.Contains(url, "mobilebanner")
}

// CheckValueInSlice will check if a value is already inside the slice.
func CheckValueInSlice(valueToCheck string, values []string) bool {
	for _, v := range values {
		if v == valueToCheck {
			return true
		}
	}
	return false
}
