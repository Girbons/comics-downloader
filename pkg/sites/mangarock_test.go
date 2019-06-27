package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestMangarockGetInfo(t *testing.T) {
	options := map[string]string{"country": "italy"}
	mr := NewMangarock(options)

	name, issueNumber := mr.GetInfo("https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049")
	assert.Equal(t, "Boruto: Naruto Next Generations", name)
	assert.Equal(t, "Vol.4 Chapter 14: Teamwork...!!", issueNumber)
}

func TestMangarockSetup(t *testing.T) {
	options := map[string]string{"country": "italy"}
	mr := NewMangarock(options)
	comic := new(core.Comic)

	comic.URLSource = "https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049"

	err := mr.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 49, len(comic.Links))
}

func TestMangarockRetrieveIssueLinks(t *testing.T) {
	options := map[string]string{"country": "italy"}
	mr := NewMangarock(options)

	issues, err := mr.RetrieveIssueLinks("https://mangarock.com/manga/mrs-serie-173467", false, false)

	assert.Nil(t, err)
	assert.Equal(t, 700, len(issues))
}

func TestMangarockRetrieveIssueLinksLastChapter(t *testing.T) {
	options := map[string]string{"country": "italy"}
	mr := NewMangarock(options)

	issues, err := mr.RetrieveIssueLinks("https://mangarock.com/manga/mrs-serie-173467", false, true)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestMangarockRetrieveLastIssueLink(t *testing.T) {
	options := map[string]string{"country": "italy"}
	mr := NewMangarock(options)

	issue, err := mr.retrieveLastIssue("https://mangarock.com/manga/mrs-serie-173467")

	assert.Nil(t, err)
	assert.Equal(t, "https://mangarock.com/manga/mrs-serie-173467/chapter/mrs-chapter-174165", issue)
}
