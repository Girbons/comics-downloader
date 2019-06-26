package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestMangatownGetInfo(t *testing.T) {
	mt := new(Mangatown)
	name, issueNumber := mt.GetInfo("http://www.mangatown.com/manga/naruto/v63/c684/")

	assert.Equal(t, "naruto", name)
	assert.Equal(t, "c684", issueNumber)
}

func TestMangatownSetup(t *testing.T) {
	mt := new(Mangatown)
	comic := new(core.Comic)
	comic.URLSource = "http://www.mangatown.com/manga/naruto/v63/c684/"

	err := mt.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 22, len(comic.Links))
}

func TestMangatownRetrieveIssueLinks(t *testing.T) {
	mt := new(Mangatown)
	issues, err := mt.RetrieveIssueLinks("https://www.mangatown.com/manga/naruto/", false, false)

	assert.Nil(t, err)
	assert.Equal(t, 748, len(issues))
}

func TestMangatownRetrieveLastIssueLink(t *testing.T) {
	mt := new(Mangatown)
	issue, err := mt.retrieveLastIssue("https://www.mangatown.com/manga/naruto/")

	assert.Nil(t, err)
	assert.Equal(t, "https://www.mangatown.com/manga/naruto/v72/c700.6/", issue)
}
