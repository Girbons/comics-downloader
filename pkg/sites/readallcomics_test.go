package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestReadallcomicsSetup(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = "https://readallcomics.com/sandman-v2-075-1989/"
	opt := &config.Options{
		URL:    "https://readallcomics.com/sandman-v2-075-1989/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	err := readallcomics.Initialize(comic)
	assert.Nil(t, err)
	// Note: This would depend on the actual page content, adjust expected count as needed
	assert.Greater(t, len(comic.Links), 0)
}

func TestReadallcomicsGetInfoSomethingIsKillingTheChildren(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readallcomics.com/something-is-killing-the-children-000-2024/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	name, issueNumber := readallcomics.GetInfo("https://readallcomics.com/something-is-killing-the-children-000-2024/")
	assert.Equal(t, "something is killing the children", name)
	assert.Equal(t, "000-2024", issueNumber)
}

func TestReadallcomicsGetInfoEmbeddedIssue(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readallcomics.com/something-is-killing-the-children-029something-is-killing-the-children-2023/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	name, issueNumber := readallcomics.GetInfo("https://readallcomics.com/something-is-killing-the-children-029something-is-killing-the-children-2023/")
	assert.Equal(t, "something is killing the children", name)
	assert.Equal(t, "029", issueNumber)
}

func TestReadallcomicsGetInfoSandmanV2(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readallcomics.com/sandman-v2-075-1989/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	name, issueNumber := readallcomics.GetInfo("https://readallcomics.com/sandman-v2-075-1989/")
	assert.Equal(t, "sandman", name)
	assert.Equal(t, "v2-075-1989", issueNumber)
}

func TestReadallcomicsGetInfoSandmanDeluxeEdition(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readallcomics.com/sandman-v2-_the_deluxe_edition-5-part-6-1989/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	name, issueNumber := readallcomics.GetInfo("https://readallcomics.com/sandman-v2-_the_deluxe_edition-5-part-6-1989/")
	assert.Equal(t, "sandman", name)
	assert.Equal(t, "v2-_the_deluxe_edition-5-part-6-1989", issueNumber)
}

func TestReadallcomicsRetrieveIssueLinks(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readallcomics.com/sandman-v2-075-1989/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	issues, err := readallcomics.RetrieveIssueLinks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
	assert.Equal(t, "https://readallcomics.com/sandman-v2-075-1989/", issues[0])
}

func TestReadallcomicsRetrieveIssueLinksFromSandmanCategory(t *testing.T) {
	opt := &config.Options{
		URL:    "http://readallcomics.com/category/sandman/",
		All:    true,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	issues, err := readallcomics.RetrieveIssueLinks()
	assert.Nil(t, err)
	assert.Greater(t, len(issues), 1)

	// Check if our test URLs are present
	expectedURLs := []string{
		"https://readallcomics.com/sandman-v2-075-1989/",
		"https://readallcomics.com/sandman-v2-_the_deluxe_edition-5-part-6-1989/",
	}

	issueSet := make(map[string]bool)
	for _, issue := range issues {
		issueSet[issue] = true
	}

	for _, expectedURL := range expectedURLs {
		assert.True(t, issueSet[expectedURL], "Expected URL %s should be found in issues", expectedURL)
	}
}

func TestReadallcomicsRetrieveIssueLinksLastFromSandman(t *testing.T) {
	opt := &config.Options{
		URL:    "http://readallcomics.com/category/sandman/",
		All:    false,
		Last:   true,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	issues, err := readallcomics.RetrieveIssueLinks()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
	// Should return the last issue from the sandman category
}

func TestReadallcomicsGetIssuesFromSandmanCategory(t *testing.T) {
	opt := &config.Options{
		URL:    "http://readallcomics.com/category/sandman/",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	readallcomics := NewReadallcomics(opt)
	issues, err := readallcomics.getIssues("http://readallcomics.com/category/sandman/")
	assert.Nil(t, err)
	assert.Greater(t, len(issues), 0)

	// Check if expected URLs are present in the results
	expectedURLs := []string{
		"https://readallcomics.com/sandman-v2-075-1989/",
		"https://readallcomics.com/sandman-v2-_the_deluxe_edition-5-part-6-1989/",
	}

	issueSet := make(map[string]bool)
	for _, issue := range issues {
		issueSet[issue] = true
	}

	for _, expectedURL := range expectedURLs {
		assert.True(t, issueSet[expectedURL], "Expected URL %s should be found in issues", expectedURL)
	}
}
