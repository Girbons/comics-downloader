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

func TestMangaReadGetInfo(t *testing.T) {
	name, issueNumber := GetInfo(URL)

	assert.Equal(t, "naruto", name)
	assert.Equal(t, "1", issueNumber)
}

func TestRetrieveMangaReaderImageLinks(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = URL
	comic.Name = "naruto"
	comic.IssueNumber = "1"
	comic.Source = SOURCE

	links, err := retrieveImageLinks(comic)

	assert.Equal(t, 53, len(links))
	assert.Nil(t, err)
}

func TestSetupMangaReader(t *testing.T) {
	comic := new(core.Comic)
	comic.Name = "naruto"
	comic.IssueNumber = "1"
	comic.URLSource = URL
	comic.Source = SOURCE

	err := Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 53, len(comic.Links))
}

func TestRetrieveIssueLinks(t *testing.T) {
	issues, err := RetrieveIssueLinks("https://www.mangareader.net/naruto", false)

	assert.Nil(t, err)
	assert.Equal(t, 700, len(issues))
}
