package mangareader

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

const (
	URL    = "https://www.mangareader.net/naruto/1/"
	SOURCE = "www.mangareader.net"
)

func TestRetrieveMangaReaderImageLinks(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = URL
	comic.Name = "naruto"
	comic.IssueNumber = "1"
	comic.Source = SOURCE

	links, err := retrieveImageLinks(comic)

	assert.Equal(t, 53, len(links))
	assert.Equal(t, nil, err)
}

func TestSetupMangaReader(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = URL
	comic.Source = SOURCE

	err := Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 53, len(comic.Links))
}
