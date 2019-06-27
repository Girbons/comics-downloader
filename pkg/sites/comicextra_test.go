package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestComicExtraSetup(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = "https://www.comicextra.com/batman-2016/chapter-58/full"

	comicextra := new(Comicextra)
	err := comicextra.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 24, len(comic.Links))
}

func TestComicExtraGetInfo(t *testing.T) {
	comicextra := new(Comicextra)
	name, issueNumber := comicextra.GetInfo("https://www.comicextra.com/batman-2016/chapter-58/full")

	assert.Equal(t, "batman-2016", name)
	assert.Equal(t, "chapter-58", issueNumber)
}

func TestComicextraRetrieveIssueLinks(t *testing.T) {
	comicextra := new(Comicextra)
	issues, err := comicextra.RetrieveIssueLinks("https://www.comicextra.com/comic/100-bullets", false, false)

	assert.Nil(t, err)
	assert.Equal(t, 100, len(issues))
}

func TestComicextraRetrieveIssueLinksLastChapter(t *testing.T) {
	comicextra := new(Comicextra)
	issues, err := comicextra.RetrieveIssueLinks("https://www.comicextra.com/comic/100-bullets", false, true)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestComicExtraRetrieveLastIssueLink(t *testing.T) {
	comicextra := new(Comicextra)
	issue, err := comicextra.retrieveLastIssue("https://www.comicextra.com/batman-2016/chapter-58/full")

	assert.Nil(t, err)
	assert.Equal(t, "https://www.comicextra.com/batman-2016/chapter-73/full", issue)
}

func TestComicExtraRetrieveLastIssueLinkNotDetail(t *testing.T) {
	comicextra := new(Comicextra)
	issue, err := comicextra.retrieveLastIssue("https://www.comicextra.com/comic/100-bullets")

	assert.Nil(t, err)
	assert.Equal(t, "https://www.comicextra.com/100-bullets/chapter-100", issue)
}
