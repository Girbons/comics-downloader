package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/mod/semver"
)

const (
	releasesEndpoint = "https://api.github.com/repos/Girbons/comics-downloader/releases?per_page=1"
	cacheTTL         = time.Minute * 15
)

// Tag specifies the current release tag. Overridden at build time via -ldflags.
var Tag = "development"

type release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

type releaseResult struct {
	checked   time.Time
	available bool
	link      string
	err       error
}

var (
	cacheMu      sync.Mutex
	cachedResult releaseResult
)

// IsNewAvailable fetches the latest release information and compares it against the current Tag.
// Results are cached for a short period to avoid repeatedly hammering the API.
func IsNewAvailable(ctx context.Context, client *http.Client) (bool, string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if client == nil {
		client = http.DefaultClient
	}

	cacheMu.Lock()
	if !cachedResult.checked.IsZero() && time.Since(cachedResult.checked) < cacheTTL {
		res := cachedResult
		cacheMu.Unlock()
		return res.available, res.link, res.err
	}
	cacheMu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesEndpoint, nil)
	if err != nil {
		return false, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		updateCache(false, "", err)
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		updateCache(false, "", err)
		return false, "", err
	}

	var releases []release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		updateCache(false, "", err)
		return false, "", err
	}
	if len(releases) == 0 {
		err = fmt.Errorf("no releases found")
		updateCache(false, "", err)
		return false, "", err
	}
	latest := releases[0]

	if !semver.IsValid(Tag) || !semver.IsValid(latest.TagName) {
		err = fmt.Errorf("invalid semver tag (current=%q latest=%q)", Tag, latest.TagName)
		updateCache(false, "", err)
		return false, "", err
	}

	if semver.Compare(Tag, latest.TagName) < 0 {
		updateCache(true, latest.HTMLURL, nil)
		return true, latest.HTMLURL, nil
	}

	updateCache(false, "", nil)
	return false, "", nil
}

func updateCache(available bool, link string, err error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cachedResult = releaseResult{
		checked:   time.Now(),
		available: available,
		link:      link,
		err:       err,
	}
}

// ResetCache clears the cached release lookup result. Intended for testing.
func ResetCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cachedResult = releaseResult{}
}
