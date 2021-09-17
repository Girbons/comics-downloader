package util

import (
	"net/url"
	"strings"
)

// TrimAndSplitURL trim tailing "/" and split url
func TrimAndSplitURL(u string) []string {
	u = strings.TrimSuffix(u, "/")
	return strings.Split(u, "/")
}

// URLSource retrieves the hostname from given url.
func URLSource(u string) (string, error) {
	parsedURL, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	return parsedURL.Hostname(), nil
}

// IsURLValid exclude those url containing `.gif` and `logo`.
func IsURLValid(url string) bool {
	invalidValues := []string{".gif", "logo", "mobilebanner", "wp-content"}
	check := true

	for _, v := range invalidValues {
		if strings.Contains(url, v) {
			check = false
			break
		}
	}

	if check {
		return strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https")
	}

	return check
}

// IsValueInSlice checks if a value is already in a slice.
func IsValueInSlice(valueToCheck string, values []string) bool {
	for _, v := range values {
		if v == valueToCheck {
			return true
		}
	}
	return false
}

// Parse escapes characters
func Parse(s string) string {
	replacer := strings.NewReplacer(
		".", " ",
		"/", "_",
		"[", "",
		"]", "",
		":", "",
		";", "",
		"!", "",
		"?", "",
	)

	return strings.Trim(replacer.Replace(s), " ")
}
