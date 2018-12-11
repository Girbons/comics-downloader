package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteLoaderMangarock(t *testing.T) {
	url := "https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049"
	result := LoadComicFromSource("mangarock.com", url)

	assert.Equal(t, "mangarock.com", result.Source)
	assert.Equal(t, url, result.URLSource)
	assert.Equal(t, "Boruto: Naruto Next Generations", result.Name)
	assert.Equal(t, "Vol.TBD Chapter 14: Teamwork...!!", result.IssueNumber)
	assert.Equal(t, 49, len(result.Links))
}

func TestSiteLoaderComicExtra(t *testing.T) {
	url := "https://www.comicextra.com/daredevil-2016/chapter-600/full"
	result := LoadComicFromSource("www.comicextra.com", url)

	assert.Equal(t, "www.comicextra.com", result.Source)
	assert.Equal(t, url, result.URLSource)
	assert.Equal(t, "daredevil-2016", result.Name)
	assert.Equal(t, "chapter-600", result.IssueNumber)
	assert.Equal(t, 45, len(result.Links))
}

func TestSiteLoaderMangahere(t *testing.T) {
	url := "http://www.mangahere.cc/manga/shingeki_no_kyojin_before_the_fall/c048/"
	result := LoadComicFromSource("www.mangahere.cc", url)
	assert.Equal(t, "www.mangahere.cc", result.Source)
	assert.Equal(t, url, result.URLSource)
	assert.Equal(t, "shingeki_no_kyojin_before_the_fall", result.Name)
	assert.Equal(t, "c048", result.IssueNumber)
	// TODO must be fixed because links are 65...
	assert.Equal(t, 130, len(result.Links))
}

func TestLoaderUnknownSource(t *testing.T) {
	url := "http://example.com"
	result := LoadComicFromSource("example.com", url)

	assert.Equal(t, "example.com", result.Source)
}
