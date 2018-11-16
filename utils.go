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
