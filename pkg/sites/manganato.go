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

func init() {
	SupportedSites["manganato"] = SupportedSite{
		IsEnabled: true,
		Loader:    func(opts *config.Options) BaseSite { return NewManganato(opts) },
	}
}

// GetInfo extracts the basic info from the given url.
func (m *Manganato) GetInfo(url string) (string, string, error) {
	name, issueNumber, err := MangaKakalotGetInfo(m.options, "manganato.com", url)
	if err != nil {

		return "", "", err
	}

	return name, issueNumber, nil
}

// Initialize loads links and metadata from manganato
func (m *Manganato) Initialize(comic *core.ComicIssue) error {
	return MangaKakalotInitialize(m.options, comic)
}

// RetrieveIssueLinks retrieve the issue links for the given comic.
func (m *Manganato) RetrieveIssueLinks() ([]string, error) {
	return MangaKakalotRetrieveIssueLinks(m.options, "manganato.com", m.options.URL)
}
