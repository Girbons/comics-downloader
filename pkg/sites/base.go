package sites

import "github.com/Girbons/comics-downloader/pkg/core"

// SupportedSites are the supported sites.
var SupportedSites = []string{
	"www.comicextra.com",
	"www.mangareader.net",
	"www.mangatown.com",
	"www.mangahere.cc",
	"mangadex.cc",
	"mangadex.org",
}

// DisabledSites are the sites that are currently disabled.
var DisabledSites = []string{
	"www.comicextra.com",
}

// BaseSite specifies an implementation of a Site which allows
// to retrieve a manga/comic basics info and imges links
type BaseSite interface {
	// Initialize will initialize the comic struct with the images link
	Initialize(comic *core.Comic) error

	// GetInfo will return the comic name and issue number
	GetInfo(url string) (string, string)

	// RetrieveIssueLinks will return the images links of a comic
	RetrieveIssueLinks(url string, all, last bool) ([]string, error)
}
