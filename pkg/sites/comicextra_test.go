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
	comic.URLSource = "https://www.comicextra.com/batman-2016/chapter-58/full"

	opt :=
		&config.Options{
			URL:    "https://www.comicextra.com/batman-2016/chapter-58/full",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	err := comicextra.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 24, len(comic.Links))
}

func TestComicExtraGetInfo(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://www.comicextra.com/comic/100-bullets",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	name, issueNumber := comicextra.GetInfo("https://www.comicextra.com/batman-2016/chapter-58/full")

	assert.Equal(t, "batman-2016", name)
	assert.Equal(t, "chapter-58", issueNumber)
}

func TestComicextraRetrieveIssueLinks(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://www.comicextra.com/comic/100-bullets",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 100, len(issues))
}

func TestComicextraRetrieveIssueLinksURLWithPage(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://www.comicextra.com/comic/100-bullets/2",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 100, len(issues))
}

func TestComicextraRetrieveIssueLinksInASinglePage(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://www.comicextra.com/comic/captain-marvel-2016",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}

	comicextra := NewComicextra(opt)
	issues, err := comicextra.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 10, len(issues))
}

func TestComicextraRetrieveIssueLinksLastChapter(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://www.comicextra.com/comic/100-bullets",
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
	issue, err := comicextra.retrieveLastIssue("https://www.comicextra.com/comic/100-bullets")

	assert.Nil(t, err)
	assert.Equal(t, "https://www.comicextra.com/100-bullets/chapter-100", issue)
}

func TestComicExtraRetrieveLastIssueLinkNotDetail(t *testing.T) {
	comicextra := new(Comicextra)
	issue, err := comicextra.retrieveLastIssue("https://www.comicextra.com/comic/100-bullets")

	assert.Nil(t, err)
	assert.Equal(t, "https://www.comicextra.com/100-bullets/chapter-100", issue)
}
