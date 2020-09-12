package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestReadComicOnlineSetup(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = "https://readcomiconline.to/Comic/Batman-2016/Issue-58?id=143175"

	readComicOnline := new(ReadComicOnline)
	err := readComicOnline.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 24, len(comic.Links))
}

func TestReadComicOnlineGetInfo(t *testing.T) {
	readComicOnline := new(ReadComicOnline)
	name, issueNumber := readComicOnline.GetInfo("https://readcomiconline.to/Comic/Batman-2016/Issue-58?id=143175")

	assert.Equal(t, "Batman-2016", name)
	assert.Equal(t, "58", issueNumber)
}

func TestReadComicOnlineRetrieveIssueLinks(t *testing.T) {
	readComicOnline := new(ReadComicOnline)
	issues, err := readComicOnline.RetrieveIssueLinks("https://readcomiconline.to/Comic/100-Bullets", false, false)

	assert.Nil(t, err)
	assert.Equal(t, 100, len(issues))
}

func TestReadComicOnlineRetrieveIssueLinksLastChapter(t *testing.T) {
	readComicOnline := new(ReadComicOnline)
	issues, err := readComicOnline.RetrieveIssueLinks("https://readcomiconline.to/Comic/100-Bullets", false, true)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestReadComicOnlineRetrieveLastIssueLink(t *testing.T) {
	readComicOnline := new(ReadComicOnline)
	issue, err := readComicOnline.retrieveLastIssue("https://readcomiconline.to/Comic/100-Bullets")

	assert.Nil(t, err)
	assert.Equal(t, "https://readcomiconline.to/Comic/100-Bullets/Issue-100-2", issue)
}
