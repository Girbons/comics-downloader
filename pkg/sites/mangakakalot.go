package sites

import (
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
)

type MangaKakalot struct {
	options *config.Options
}

// NewMangaKakalot returns a new MangaKakalot instance.
func NewMangaKakalot(options *config.Options) *MangaKakalot {
	return &MangaKakalot{
		options: options,
	}
}

// GetInfo extracts the basic info from the given url.
func (m *MangaKakalot) GetInfo(url string) (string, string) {
	name, issueNumber, err := MangaKakalotGetInfo(m.options, "mangakakalot.com", url)
	if err != nil {
		m.options.Logger.Errorf("error getting info for url %q: %v", url, err)
		return "", ""
	}

	return name, issueNumber
}

// Initialize loads links and metadata from mangakakalot
func (m *MangaKakalot) Initialize(comic *core.ComicIssue) error {
	return MangaKakalotInitialize(m.options, comic)
}

// RetrieveIssueLinks retrieve the issue links for the given comic.
func (m *MangaKakalot) RetrieveIssueLinks() ([]string, error) {
	return MangaKakalotRetrieveIssueLinks(m.options, "mangakakalot.com", m.options.URL)
}
