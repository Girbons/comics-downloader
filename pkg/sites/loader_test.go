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
	url := "https://comicextra.me/batman-unseen/issue-5/full"
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
	assert.Equal(t, 1, len(collection))

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
	url := "https://comicextra.net/batman-unseen/issue-5/full"
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

func TestVolumeIssuesRange(t *testing.T) {
	tt := []struct {
		name        string
		input       string
		start       float64
		end         float64
		returnValue bool
	}{
		// Volume 4, Issue 78-99 range tests (user specifies: -range=4.78-4.99)
		{"v4-078 in range", "v4-078-2016", 4.78, 4.99, false},
		{"v4-099 in range", "v4-099-2016", 4.78, 4.99, false},
		{"v4-077 out of range (too low)", "v4-077-2016", 4.78, 4.99, true},
		{"v4-100 out of range (too high)", "v4-100-2016", 4.78, 4.99, true},
		{"v3-078 wrong volume", "v3-078-2016", 4.78, 4.99, true},
		{"v5-078 wrong volume", "v5-078-2016", 4.78, 4.99, true},

		// Volume 2, Issue 1-50 range tests (user specifies: -range=2.01-2.50)
		{"v2-001 in range", "v2-001-1989", 2.01, 2.50, false},
		{"v2-025 in range", "v2-025-1989", 2.01, 2.50, false},
		{"v2-050 in range", "v2-050-1989", 2.01, 2.50, false},
		{"v2-051 out of range", "v2-051-1989", 2.01, 2.50, true},

		// Edge cases with different formats
		{"v4-078 without year", "v4-078", 4.78, 4.99, false},
		{"v4_078 with underscore", "v4_078", 4.78, 4.99, false},
		{"v10-005 two-digit volume", "v10-005", 10.05, 10.10, false},

		// Backwards compatibility with simple numeric issues
		{"078 simple format", "078", 78, 99, false},
		{"100 simple format out of range", "100", 78, 99, true},
		{"issue-1 with prefix", "issue-1", 1, 3, false},
		{"issue-5 with prefix out of range", "issue-5", 1, 3, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.returnValue, notInIssuesRange(tc.input, tc.start, tc.end),
				"Issue %s with range %.2f-%.2f", tc.input, tc.start, tc.end)
		})
	}
}

func TestExtractIssueNumberForRange(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected float64
	}{
		// Volume and issue format
		{"v4-078-2016", "v4-078-2016", 4.78},
		{"v4-099-2016", "v4-099-2016", 4.99},
		{"v2-075-1989", "v2-075-1989", 2.75},
		{"v4-078 no year", "v4-078", 4.78},
		{"v4_078 underscore", "v4_078", 4.78},
		{"v10-005 two-digit volume", "v10-005", 10.05},

		// Simple numeric format
		{"078", "078", 78},
		{"99", "99", 99},
		{"1", "1", 1},

		// Decimal format
		{"20.5", "20.5", 20.5},
		{"3.14", "3.14", 3.14},

		// With prefixes
		{"issue-1", "issue-1", 1},
		{"issue-123", "issue-123", 123},

		// Edge cases
		{"empty", "", 0},
		{"no numbers", "abc", 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, extractIssueNumberForRange(tc.input),
				"Issue %s should extract to %.2f", tc.input, tc.expected)
		})
	}
}
