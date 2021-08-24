package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestReadComicOnlineSetup(t *testing.T) {

	comic := new(core.Comic)
	comic.URLSource = "https://readcomiconline.li/Comic/Batman-2016/Issue-58?id=143175"

	opt :=
		&config.Options{
			URL:    "https://readcomiconline.li/Comic/Batman-2016/Issue-58?id=143175",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	readComicOnline := NewReadComiconline(opt)
	err := readComicOnline.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 24, len(comic.Links))
}

func TestReadComicOnlineGetInfo(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://readcomiconline.li/Comic/Batman-2016/Issue-58?id=143175",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	readComicOnline := NewReadComiconline(opt)
	name, issueNumber := readComicOnline.GetInfo("https://readcomiconline.li/Comic/Batman-2016/Issue-58?id=143175")

	assert.Equal(t, "Batman-2016", name)
	assert.Equal(t, "58", issueNumber)
}

func TestReadComicOnlineRetrieveIssueLinks(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://readcomiconline.li/Comic/100-Bullets",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	readComicOnline := NewReadComiconline(opt)
	issues, err := readComicOnline.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 100, len(issues))
}

func TestReadComicOnlineRetrieveIssueLinksLastChapter(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://readcomiconline.li/Comic/100-Bullets",
			All:    false,
			Last:   true,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	readComicOnline := NewReadComiconline(opt)
	issues, err := readComicOnline.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestReadComicOnlineRetrieveLastIssueLink(t *testing.T) {
	opt :=
		&config.Options{
			URL:    "https://readcomiconline.li/Comic/100-Bullets",
			All:    false,
			Last:   true,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	readComicOnline := NewReadComiconline(opt)
	issue, err := readComicOnline.retrieveLastIssue("https://readcomiconline.li/Comic/100-Bullets")

	assert.Nil(t, err)
	assert.Equal(t, "https://readcomiconline.li/Comic/100-Bullets/Issue-100-2", issue)
}
