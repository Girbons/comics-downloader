package sites

import (
	"github.com/Girbons/comics-downloader/pkg/core"
)

type BaseSite interface {
	Initialize(comic *core.Comic) error
	GetInfo(url string, options map[string]string) (string, string)
	RetrieveIssueLinks(url string, all bool, options map[string]string) ([]string, error)
}

func Initialize(b BaseSite, comic *core.Comic) error {
	return b.Initialize(comic)
}

func GetInfo(b BaseSite, url string, options map[string]string) (string, string) {
	return b.GetInfo(url, options)
}

func RetrieveIssueLinks(b BaseSite, url string, all bool, options map[string]string) ([]string, error) {
	return b.RetrieveIssueLinks(url, all, options)
}

type SiteLoader struct {
	Source BaseSite
}
