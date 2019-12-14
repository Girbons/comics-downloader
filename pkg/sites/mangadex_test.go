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
	urls, err := md.RetrieveIssueLinks("https://mangadex.org/chapter/155061/", false, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(urls))
}

func TestMangadexRetrieveIssueLinksAllChapter(t *testing.T) {
	md := NewMangadex("gb")
	urls, err := md.RetrieveIssueLinks("https://mangadex.org/title/5/naruto/", true, false)
	assert.Nil(t, err)
	assert.Len(t, urls, 569)
}

func TestMangadexRetrieveIssueLinksLastChapter(t *testing.T) {
	md := NewMangadex("gb")
	urls, err := md.RetrieveIssueLinks("https://mangadex.org/title/5/naruto/", false, true)
	assert.Nil(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, "https://mangadex.org/chapter/670438", urls[0])
}

func TestMangadexUnsupportedURL(t *testing.T) {
	md := NewMangadex("")
	_, err := md.RetrieveIssueLinks("https://mangadex.org/", false, false)
	assert.EqualError(t, err, "URL not supported")
	_, err = md.RetrieveIssueLinks("https://mangadex.org/test/0/", false, false)
	assert.EqualError(t, err, "URL not supported")
}

func TestMangadexNoManga(t *testing.T) {
	md := NewMangadex("")
	_, err := md.RetrieveIssueLinks("https://mangadex.org/title/0/", false, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Manga ID does not exist")
}

func TestMangadexNoChapters(t *testing.T) {
	md := NewMangadex("xyz")
	_, err := md.RetrieveIssueLinks("https://mangadex.org/title/5/naruto/", true, false)
	assert.EqualError(t, err, "no chapters found")
}
