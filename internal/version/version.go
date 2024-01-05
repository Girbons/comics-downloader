package version

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/mod/semver"
)

// Tag specifies the current release tag.
// It needs to be manually updated.
const Tag = "v0.33.8"

// IsNewAvailable will fetch the latest project releases
// and will compare the latest release Tag against the current Tag.
func IsNewAvailable() (bool, string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(ctx, "Girbons", "comics-downloader", nil)

	if err != nil {
		return false, "", err
	}

	// Compare returns an integer comparing two versions
	// according to semantic version precedence.
	result := semver.Compare(Tag, *releases[0].TagName)

	// -1 if v < w
	if result == -1 {
		return true, *releases[0].HTMLURL, err
	}

	return false, "", err
}
