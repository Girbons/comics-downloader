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
	return MangaKakalotGetInfo("mangakakalot.com", url)
}

// Initialize loads links and metadata from mangakakalot
func (m *MangaKakalot) Initialize(comic *core.Comic) error {
	return MangaKakalotInitialize(comic)
}

// RetrieveIssueLinks retrieve the issue links for the given comic.
func (m *MangaKakalot) RetrieveIssueLinks() ([]string, error) {
	return MangaKakalotRetrieveIssueLinks("mangakakalot.com", m.options.URL)
}
