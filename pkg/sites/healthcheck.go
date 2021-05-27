package sites

import (
	"fmt"
	"net/http"
)

// Healthcheck checks that a site is available.
func Healthcheck(url string) (bool, string) {
	res, err := http.Get(url)

	if err != nil {
		return false, err.Error()
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return true, ""
	}

	return false, fmt.Sprintf("health check failed, response status %d", res.StatusCode)
}
