package sites

import (
	"github.com/Girbons/comics-downloader/pkg/core"
)

type BaseSite interface {
	Initialize(comic *core.Comic) error
	GetInfo(url string) (string, string)
	RetrieveIssueLinks(url string, all bool) ([]string, error)
}

func Initialize(b BaseSite, comic *core.Comic) error {
	return b.Initialize(comic)
}

func GetInfo(b BaseSite, url string) (string, string) {
	return b.GetInfo(url)
}

func RetrieveIssueLinks(b BaseSite, url string, all bool) ([]string, error) {
	return b.RetrieveIssueLinks(url, all)
}

type SiteLoader struct {
	Source BaseSite
}
