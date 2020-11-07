package sites

import "github.com/Girbons/comics-downloader/pkg/core"

// BaseSite specifies an implementation of a Site which allows
// to retrieve a manga/comic basics info and imges links
type BaseSite interface {
	// Initialize will initialize the comic struct with the images link
	Initialize(comic *core.Comic) error

	// GetInfo will return the comic name and issue number
	GetInfo(url string) (string, string)

	// RetrieveIssueLinks will return the images links of a comic
	RetrieveIssueLinks() ([]string, error)
}
