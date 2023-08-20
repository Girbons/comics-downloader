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
		URL:          url,
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

func TestCustomComicName(t *testing.T) {
	url := "https://www.mangatown.com/manga/naruto/v63/c693/"
	outputFolder := filepath.Dir(os.Args[0])

	options := &config.Options{
		All:             false,
		Last:            false,
		ImagesOnly:      false,
		Source:          "www.mangatown.com",
		URL:             url,
		Format:          "pdf",
		ImagesFormat:    "png",
		CustomComicName: "Naruto",
		OutputFolder:    outputFolder,
	}

	collection, err := LoadComicFromSource(options)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 1)

	comic := collection[0]

	assert.Equal(t, "www.mangatown.com", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "Naruto", comic.Name)
	assert.Equal(t, "c693", comic.IssueNumber)
	assert.Equal(t, 20, len(comic.Links))
}

//func TestSiteLoaderMangareader(t *testing.T) {
//url := "https://www.mangareader.net/naruto/700"
//outputFolder := filepath.Dir(os.Args[0])

//options := &config.Options{
//All:          false,
//Last:         false,
//ImagesOnly:   false,
//Source:       "www.mangareader.net",
//Url:          url,
//Format:       "pdf",
//ImagesFormat: "png",
//OutputFolder: outputFolder,
//}

//collection, err := LoadComicFromSource(options)

//assert.Nil(t, err)
//assert.Equal(t, len(collection), 1)

//comic := collection[0]

//assert.Equal(t, "www.mangareader.net", comic.Source)
//assert.Equal(t, url, comic.URLSource)
//assert.Equal(t, "naruto", comic.Name)
//assert.Equal(t, "700", comic.IssueNumber)
//assert.Equal(t, 23, len(comic.Links))
//}

func TestSiteLoaderComicExtra(t *testing.T) {
  url := "https://comicextra.net/comic/batman-unseen/issue-5/full"
	outputFolder := filepath.Dir(os.Args[0])
	options := &config.Options{
		All:          false,
		Last:         false,
		ImagesOnly:   false,
		Source:       "comicextra.net",
		URL:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
	}
	collection, err := LoadComicFromSource(options)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(collection))

	comic := collection[0]

	assert.Equal(t, "comicextra.net", comic.Source)
	assert.Equal(t, url, comic.URLSource)
	assert.Equal(t, "batman-unseen", comic.Name)
	assert.Equal(t, "issue-5", comic.IssueNumber)
	assert.Equal(t, 23, len(comic.Links))
}

func TestLoaderUnknownSource(t *testing.T) {
	url := "http://example.com"
	outputFolder := filepath.Dir(os.Args[0])

	options := &config.Options{
		All:          false,
		Last:         false,
		ImagesOnly:   false,
		Source:       "example.com",
		URL:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
	}

	collection, err := LoadComicFromSource(options)

	if assert.NotNil(t, err) {
		assert.Equal(t, fmt.Errorf("source unknown"), err)
	}
	assert.Equal(t, len(collection), 0)
}

func TestIssuesRange(t *testing.T) {
  url := "https://comicextra.net/comic/batman-unseen/issue-5/full"
	outputFolder := filepath.Dir(os.Args[0])
	options := &config.Options{
		All:          true,
		Last:         false,
		ImagesOnly:   false,
		Source:       "comicextra.net",
		URL:          url,
		Format:       "pdf",
		ImagesFormat: "png",
		OutputFolder: outputFolder,
		IssuesRange:  "1-3",
	}
	collection, err := LoadComicFromSource(options)

	assert.Nil(t, err)
	assert.Equal(t, len(collection), 3)

	issues := make([]string, 0, len(collection))
	for _, c := range collection {
		issues = append(issues, c.IssueNumber)
	}

	assert.Contains(t, issues, "issue-1")
	assert.Contains(t, issues, "issue-2")
	assert.Contains(t, issues, "issue-3")
}

func TestFloatIssuesRange(t *testing.T) {
	tt := []struct {
		input       string
		start       float64
		end         float64
		returnValue bool
	}{
		{"1", 1, 1, false},
		{"19", 20, 21, true},
		{"20", 20, 21, false},
		{"20.5", 20, 21, false},
		{"21", 20, 21, false},
		{"22", 20, 21, true},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, notInIssuesRange(tc.input, tc.start, tc.end), tc.returnValue)
		})
	}
}
