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
	comic.URLSource = "https://comicextra.me/batman-unseen/issue-5/full"

	opt :=
		&config.Options{
			URL:    "https://comicextra.me/batman-unseen/issue-5/full",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	err := comicextra.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 23, len(comic.Links))
}

func TestComicExtraGetInfo(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://comicextra.me/batman-unseen/issue-5/full",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	name, issueNumber := comicextra.GetInfo("https://comicextra.me/batman-unseen/issue-5/full")

	assert.Equal(t, "batman-unseen", name)
	assert.Equal(t, "issue-5", issueNumber)
}

func TestComicextraRetrieveIssueLinks(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://comicextra.me/batman-unseen/issue-5/full",
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

func TestComicextraRetrieveIssueLinksURLWithPage(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://comicextra.me/batman-unseen/full",
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

func TestComicextraRetrieveIssueLinksInASinglePage(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://comicextra.me/batman-unseen/issue-4/full",
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
			URL:    "https://comicextra.me/batman-unseen/issue-4/full",
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
	issue, err := comicextra.retrieveLastIssue("https://comicextra.me/batman-unseen/issue-1/full")

	assert.Nil(t, err)
	assert.Equal(t, "https://comicextra.me/batman-unseen/issue-5/full", issue)
}

func TestComicExtraRetrieveLastIssueLinkNotDetail(t *testing.T) {
	comicextra := new(Comicextra)
	issue, err := comicextra.retrieveLastIssue("https://comicextra.me/batman-unseen/issue-1/full")

	assert.Nil(t, err)
	assert.Equal(t, "https://comicextra.me/batman-unseen/issue-5/full", issue)
}
