package main

import (
	"net/url"
	"strings"
)

func SplitUrl(u string) []string {
	return strings.Split(u, "/")
}

func UrlSource(u string) (string, error) {
	parsedUrl, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	return parsedUrl.Hostname(), nil
}
