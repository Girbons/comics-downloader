package sites

import (
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
)

type Manganato struct {
	options *config.Options
}

// NewManganato returns a new Manganato instance.
func NewManganato(options *config.Options) *Manganato {
	return &Manganato{
		options: options,
	}
}

// GetInfo extracts the basic info from the given url.
func (m *Manganato) GetInfo(url string) (string, string) {
	return MangaKakalotGetInfo("manganato.com", url)
}

// Initialize loads links and metadata from manganato
func (m *Manganato) Initialize(comic *core.Comic) error {
	return MangaKakalotInitialize(comic)
}

// RetrieveIssueLinks retrieve the issue links for the given comic.
func (m *Manganato) RetrieveIssueLinks() ([]string, error) {
	return MangaKakalotRetrieveIssueLinks("manganato.com", m.options.URL)
}
