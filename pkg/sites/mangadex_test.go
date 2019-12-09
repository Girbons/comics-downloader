package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestMangadexGetInfo(t *testing.T) {
	md := NewMangadex("")

	name, issueNumber := md.GetInfo("https://mangadex.org/chapter/155061/1")
	assert.Equal(t, "Naruto", name)
	assert.Equal(t, "Vol 60 Chapter 575, A Will of Stone", issueNumber)
}

func TestMangadexSetup(t *testing.T) {
	md := NewMangadex("")
	comic := new(core.Comic)

	comic.URLSource = "https://mangadex.org/chapter/155061/1"

	err := md.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 14, len(comic.Links))
}

func TestMangadexRetrieveIssueLinks(t *testing.T) {
	md := NewMangadex("")

	issues, err := md.RetrieveIssueLinks("https://mangadex.org/chapter/155061/1", false, false)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}
