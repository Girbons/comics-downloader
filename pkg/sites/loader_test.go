package sites

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestSiteLoaderMangatown(t *testing.T) {
	url := "https://www.mangatown.com/manga/naruto/v63/c693/"
	outputFolder := filepath.Dir(os.Args[0])

	options := &config.Options{
		All:          false,
		Last:         false,
		ImagesOnly:   false,
		Source:       "www.mangatown.com",
		Url:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
	}

	collection, err := LoadComicFromSource(options)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 1)

	comic := collection[0]

	assert.Equal(t, "www.mangatown.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "naruto", comic.Name)
	assert.Equal(t, "c693", comic.IssueNumber)
	assert.Equal(t, 20, len(comic.Links))
}

func TestSiteLoaderMangareader(t *testing.T) {
	url := "https://www.mangareader.net/naruto/700"
	outputFolder := filepath.Dir(os.Args[0])

	options := &config.Options{
		All:          false,
		Last:         false,
		ImagesOnly:   false,
		Source:       "www.mangareader.net",
		Url:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
	}

	collection, err := LoadComicFromSource(options)

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
	outputFolder := filepath.Dir(os.Args[0])
	options := &config.Options{
		All:          false,
		Last:         false,
		ImagesOnly:   false,
		Source:       "www.comicextra.com",
		Url:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
	}
	collection, err := LoadComicFromSource(options)

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
	outputFolder := filepath.Dir(os.Args[0])

	options := &config.Options{
		All:          false,
		Last:         false,
		ImagesOnly:   false,
		Source:       "example.com",
		Url:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
	}

	collection, err := LoadComicFromSource(options)

	if assert.NotNil(t, err) {
		assert.Equal(t, fmt.Errorf("It was not possible to determine the source"), err)
	}
	assert.Equal(t, len(collection), 0)
}
