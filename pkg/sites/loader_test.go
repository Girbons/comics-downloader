package sites

import (
	"fmt"
	"testing"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestSiteLoaderMangarock(t *testing.T) {
	url := "https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049"
	conf := new(config.ComicConfig)

	comic, err := LoadComicFromSource(conf, "mangarock.com", url, "italy", "")

	assert.Nil(t, err)
	assert.Equal(t, "mangarock.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "Boruto: Naruto Next Generations", comic.Name)
	assert.Equal(t, "Vol.4 Chapter 14: Teamwork...!!", comic.IssueNumber)
	assert.Equal(t, 49, len(comic.Links))
}

func TestSiteLoaderComicExtra(t *testing.T) {
	url := "https://www.comicextra.com/daredevil-2016/chapter-600/full"
	conf := new(config.ComicConfig)
	comic, err := LoadComicFromSource(conf, "www.comicextra.com", url, "", "")

	assert.Nil(t, err)
	assert.Equal(t, "www.comicextra.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "daredevil-2016", comic.Name)
	assert.Equal(t, "chapter-600", comic.IssueNumber)
	assert.Equal(t, 43, len(comic.Links))
}

func TestSiteLoaderMangahereIsDisabled(t *testing.T) {
	url := "http://www.mangahere.cc/manga/shingeki_no_kyojin_before_the_fall/c048/"
	conf := new(config.ComicConfig)

	_, err := LoadComicFromSource(conf, "www.mangahere.cc", url, "", "")

	assert.EqualError(t, err, "mangahere is currently disabled")
}

func TestLoaderUnknownSource(t *testing.T) {
	url := "http://example.com"
	conf := new(config.ComicConfig)

	comic, err := LoadComicFromSource(conf, "example.com", url, "", "")

	if assert.NotNil(t, err) {
		assert.Equal(t, fmt.Errorf("It was not possible to determine the source"), err)
	}
	assert.Equal(t, "example.com", comic.Source)
}
