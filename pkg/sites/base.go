package sites

import "github.com/Girbons/comics-downloader/pkg/core"

// SupportedSites are the supported sites.
var SupportedSites = []string{
	"www.comicextra.com",
	"mangarock.com",
	"www.mangareader.net",
	"www.mangatown.com",
	"www.mangahere.cc",
}

// BaseSite specifies an implementation of a Site which allows you
// to retrieve a manga/comic basics info and imges links
type BaseSite interface {
	// Initialize will initialize the comic struct with the images link
	Initialize(comic *core.Comic) error

	// GetInfo will return the comic name and issue number
	GetInfo(url string) (string, string)

	// RetrieveIssueLinks will return the images links of a comic
	RetrieveIssueLinks(url string, all, last bool) ([]string, error)
}

func Initialize(b BaseSite, comic *core.Comic) error {
	return b.Initialize(comic)
}

func GetInfo(b BaseSite, url string) (string, string) {
	return b.GetInfo(url)
}

func RetrieveIssueLinks(b BaseSite, url string, all, last bool) ([]string, error) {
	return b.RetrieveIssueLinks(url, all, last)
}
