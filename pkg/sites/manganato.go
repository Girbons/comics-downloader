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
	name, issueNumber, err := MangaKakalotGetInfo(m.options, "manganato.com", url)
	if err != nil {
		m.options.Logger.Errorf("error getting info for url %q: %v", url, err)
		return "", ""
	}

	return name, issueNumber
}

// Initialize loads links and metadata from manganato
func (m *Manganato) Initialize(comic *core.ComicIssue) error {
	return MangaKakalotInitialize(m.options, comic)
}

// RetrieveIssueLinks retrieve the issue links for the given comic.
func (m *Manganato) RetrieveIssueLinks() ([]string, error) {
	return MangaKakalotRetrieveIssueLinks(m.options, "manganato.com", m.options.URL)
}
