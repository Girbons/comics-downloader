package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestMangatownGetInfo(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "http://www.mangatown.com/manga/naruto/v63/c684/",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	mt := NewMangatown(opt)
	name, issueNumber := mt.GetInfo("http://www.mangatown.com/manga/naruto/v63/c684/")

	assert.Equal(t, "naruto", name)
	assert.Equal(t, "c684", issueNumber)
}

func TestMangatownSetup(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "http://www.mangatown.com/manga/naruto/v63/c684/",
			All:    false,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	mt := NewMangatown(opt)
	comic := new(core.Comic)
	comic.URLSource = "http://www.mangatown.com/manga/naruto/v63/c684/"

	err := mt.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 22, len(comic.Links))
}

func TestMangatownRetrieveIssueLinks(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "http://www.mangatown.com/manga/naruto/v63/c684/",
			All:    true,
			Last:   false,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	mt := NewMangatown(opt)
	issues, err := mt.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 752, len(issues))
}

func TestMangatownRetrieveIssueLinksLastChapter(t *testing.T) {
	opt :=
		&config.Options{
			Url:    "http://www.mangatown.com/manga/naruto/",
			All:    false,
			Last:   true,
			Debug:  false,
			Logger: logger.NewLogger(false, make(chan string)),
		}
	mt := NewMangatown(opt)
	issues, err := mt.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}
