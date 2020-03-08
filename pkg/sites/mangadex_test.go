package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

const testMangadexBase string = "mangadex.org"
const testMangadexURL string = "https://" + testMangadexBase + "/"

func TestMangadexGetInfo(t *testing.T) {
	md := NewMangadex("", testMangadexBase)

	name, issueNumber := md.GetInfo(testMangadexURL + "chapter/155061/1")
	assert.Equal(t, "Naruto", name)
	assert.Equal(t, "Vol 60 Chapter 575, A Will of Stone", issueNumber)
}

func TestMangadexSetup(t *testing.T) {
	md := NewMangadex("", testMangadexBase)
	comic := new(core.Comic)

	comic.URLSource = testMangadexURL + "chapter/155061/1"

	err := md.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 14, len(comic.Links))
}

func TestMangadexRetrieveIssueLinks(t *testing.T) {
	md := NewMangadex("", testMangadexBase)
	urls, err := md.RetrieveIssueLinks(testMangadexURL+"chapter/155061/", false, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(urls))
}

func TestMangadexRetrieveIssueLinksAllChapter(t *testing.T) {
	md := NewMangadex("gb", testMangadexBase)
	urls, err := md.RetrieveIssueLinks(testMangadexURL+"title/5/naruto/", true, false)
	assert.Nil(t, err)
	assert.Len(t, urls, 569)
}

func TestMangadexRetrieveIssueLinksLastChapter(t *testing.T) {
	md := NewMangadex("gb", testMangadexBase)
	urls, err := md.RetrieveIssueLinks(testMangadexURL+"title/5/naruto/", false, true)
	assert.Nil(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, testMangadexURL+"chapter/670438", urls[0])
}

func TestMangadexUnsupportedURL(t *testing.T) {
	md := NewMangadex("", testMangadexBase)
	_, err := md.RetrieveIssueLinks(testMangadexURL, false, false)
	assert.EqualError(t, err, "URL not supported")
	_, err = md.RetrieveIssueLinks(testMangadexURL+"test/0/", false, false)
	assert.EqualError(t, err, "URL not supported")
}

func TestMangadexNoManga(t *testing.T) {
	md := NewMangadex("", testMangadexBase)
	_, err := md.RetrieveIssueLinks(testMangadexURL+"title/0/", false, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not get manga 0")
}

func TestMangadexNoChapters(t *testing.T) {
	md := NewMangadex("xyz", testMangadexBase)
	_, err := md.RetrieveIssueLinks(testMangadexURL+"title/5/naruto/", true, false)
	assert.EqualError(t, err, "no chapters found")
}
