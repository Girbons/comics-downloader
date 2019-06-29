package sites

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteLoaderMangatown(t *testing.T) {
	url := "https://www.mangatown.com/manga/naruto/v63/c693/"
	collection, err := LoadComicFromSource("www.mangatown.com", url, "", "pdf", false, false)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 1)

	comic := collection[0]

	assert.Equal(t, "www.mangatown.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "naruto", comic.Name)
	assert.Equal(t, "c693", comic.IssueNumber)
	assert.Equal(t, 20, len(comic.Links))
}

func TestSiteLoaderMangarock(t *testing.T) {
	url := "https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049"
	collection, err := LoadComicFromSource("mangarock.com", url, "italy", "pdf", false, false)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 1)

	comic := collection[0]

	assert.Equal(t, "mangarock.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "Boruto: Naruto Next Generations", comic.Name)
	assert.Equal(t, "Vol.4 Chapter 14: Teamwork...!!", comic.IssueNumber)
	assert.Equal(t, 49, len(comic.Links))
}

func TestSiteLoaderMangareader(t *testing.T) {
	url := "https://www.mangareader.net/naruto/700"
	collection, err := LoadComicFromSource("www.mangareader.net", url, "", "pdf", false, false)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 1)

	comic := collection[0]

	assert.Equal(t, "www.mangareader.net", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "naruto", comic.Name)
	assert.Equal(t, "700", comic.IssueNumber)
	assert.Equal(t, 23, len(comic.Links))
}

func TestSiteLoaderComicExtra(t *testing.T) {
	url := "https://www.comicextra.com/daredevil-2016/chapter-600/full"
	collection, err := LoadComicFromSource("www.comicextra.com", url, "", "pdf", false, false)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 1)

	comic := collection[0]

	assert.Equal(t, "www.comicextra.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "daredevil-2016", comic.Name)
	assert.Equal(t, "chapter-600", comic.IssueNumber)
	assert.Equal(t, 43, len(comic.Links))
}

func TestLoaderUnknownSource(t *testing.T) {
	url := "http://example.com"

	collection, err := LoadComicFromSource("example.com", url, "", "pdf", false, false)

	if assert.NotNil(t, err) {
		assert.Equal(t, fmt.Errorf("It was not possible to determine the source"), err)
	}
	assert.Equal(t, len(collection), 0)
}
