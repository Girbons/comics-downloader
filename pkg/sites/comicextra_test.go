package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestComicExtraSetup(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full"

	opt :=
		&config.Options{
			URL:    "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	err := comicextra.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 25, len(comic.Links))
}

func TestComicExtraGetInfo(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	name, issueNumber := comicextra.GetInfo("https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full")

	assert.Equal(t, "injustice-gods-among-us-year-four", name)
	assert.Equal(t, "issue-24", issueNumber)
}

func TestComicextraRetrieveIssueLinks(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://ww1.comicextra.com/comic/injustice-gods-among-us-year-four",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 25, len(issues))
}

func TestComicextraRetrieveIssueLinksURLWithPage(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://ww1.comicextra.com/comic/injustice-gods-among-us-year-four",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 25, len(issues))
}

func TestComicextraRetrieveIssueLinksInASinglePage(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestComicextraRetrieveIssueLinksLastChapter(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full",
			All:    false,
			Last:   true,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestComicExtraRetrieveLastIssueLink(t *testing.T) {
	comicextra := new(Comicextra)
	issue, err := comicextra.retrieveLastIssue("https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full")

	assert.Nil(t, err)
	assert.Equal(t, "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-annual-1/full", issue)
}

func TestComicExtraRetrieveLastIssueLinkNotDetail(t *testing.T) {
	comicextra := new(Comicextra)
	issue, err := comicextra.retrieveLastIssue("https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-24/full")

	assert.Nil(t, err)
	assert.Equal(t, "https://ww1.comicextra.com/injustice-gods-among-us-year-four/issue-annual-1/full", issue)
}
