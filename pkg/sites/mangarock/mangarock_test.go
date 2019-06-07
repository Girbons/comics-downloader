package mangarock

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestMangarockGetInfo(t *testing.T) {
	options := map[string]string{"country": "italy"}

	name, issueNumber := GetInfo("https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049", options)
	assert.Equal(t, "Boruto: Naruto Next Generations", name)
	assert.Equal(t, "Vol.4 Chapter 14: Teamwork...!!", issueNumber)
}

func TestMangarockSetup(t *testing.T) {
	comic := new(core.Comic)
	options := map[string]string{"country": "italy"}
	comic.Options = options
	comic.URLSource = "https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049"

	err := Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 49, len(comic.Links))
}

func TestRetrieveIssueLinks(t *testing.T) {
	options := map[string]string{"country": "italy"}
	issues, err := RetrieveIssueLinks("https://mangarock.com/manga/mrs-serie-173467", false, options)

	assert.Nil(t, err)
	assert.Equal(t, 700, len(issues))
}
