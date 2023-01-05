package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestMangaKakalotGetInfo(t *testing.T) {
	opt := &config.Options{
		URL:    "https://mangakakalot.com/read-wa8ap158524529412",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewMangaKakalot(opt)
	// mangakakalot.com
	name, issueNumber := mk.GetInfo("https://mangakakalot.com/chapter/evergreen/chapter_2")
	assert.Equal(t, "A Case", name)
	assert.Equal(t, "2", issueNumber)
}

func TestMangaKakalotSetup(t *testing.T) {
	opt := &config.Options{
		URL:    "https://mangakakalot.com/read-wa8ap158524529412",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewMangaKakalot(opt)
	comic := new(core.Comic)
	comic.URLSource = "https://mangakakalot.com/chapter/evergreen/chapter_2"

	err := mk.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 40, len(comic.Links))
}

func TestMangaKakalotRetrieveIssueLinks(t *testing.T) {
	opt := &config.Options{
		URL:    "https://mangakakalot.com/read-wa8ap158524529412",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewMangaKakalot(opt)
	links, err := mk.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 46, len(links))
}
