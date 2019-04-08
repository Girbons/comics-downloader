package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectComicExtra(t *testing.T) {
	source, check := DetectComic("http://www.comicextra.com/daredevil-2016/chapter-600/full")

	assert.True(t, check)
	assert.Equal(t, "www.comicextra.com", source)
}

//func TestDetectMangaHere(t *testing.T) {
//source, check := DetectComic("http://www.mangahere.cc/manga/shingeki_no_kyojin_before_the_fall/c048/")

//assert.False(t, check)
//assert.Equal(t, "www.mangahere.cc", source)
//}

func TestDetectMangaRock(t *testing.T) {
	source, check := DetectComic("https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049")

	assert.True(t, check)
	assert.Equal(t, "mangarock.com", source)
}

func TestUnsupportedSource(t *testing.T) {
	source, check := DetectComic("http://example.com")

	assert.False(t, check)
	assert.Equal(t, "", source)
}
