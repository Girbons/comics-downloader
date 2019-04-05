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
	comic.SetURLSource(URL)
	comic.SetName("naruto")
	comic.SetIssueNumber("1")
	comic.SetSource(SOURCE)

	links, err := retrieveImageLinks(comic)

	assert.Equal(t, 53, len(links))
	assert.Equal(t, nil, err)
}

func TestSetupMangaReader(t *testing.T) {
	comic := new(core.Comic)
	comic.SetURLSource(URL)
	comic.SetSource(SOURCE)

	err := Initialize(comic)
	comic.SetSource("www.mangareader.net")

	assert.Nil(t, err)
	assert.Equal(t, 53, len(comic.Links))
}
