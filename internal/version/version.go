package version

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/mod/semver"
)

// Tag specifies the current release tag.
// It needs to be manually updated.
const Tag = "v0.23.0"

// IsNewAvailable() will fetch the latest project releases
// and will compare the latest release Tag against the current Tag.
func IsNewAvailable() (bool, string) {
	ctx := context.Background()
	client := github.NewClient(nil)
	res, _, _ := client.Repositories.ListReleases(ctx, "Girbons", "comics-downloader", nil)
	// Compare returns an integer comparing two versions
	// according to semantic version precedence.
	result := semver.Compare(Tag, *res[0].TagName)

	// -1 if v < w
	if result == -1 {
		return true, *res[0].HTMLURL
	}

	return false, ""
}
