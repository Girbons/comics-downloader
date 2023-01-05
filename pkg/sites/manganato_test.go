package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestManganatoGetInfo(t *testing.T) {
	opt := &config.Options{
		URL:    "https://chapmanganato.com/manga-ng952689",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mg := NewManganato(opt)
	// readmanganato.com
	name, issueNumber := mg.GetInfo("https://chapmanganato.com/manga-ng952689/chapter-700.5")
	assert.Equal(t, "Uzumaki Naruto", name)
	assert.Equal(t, "700.5", issueNumber)
}

func TestManganatoSetup(t *testing.T) {
	opt := &config.Options{
		URL:    "https://chapmanganato.com/manga-ng952689",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewManganato(opt)
	comic := new(core.Comic)
	comic.URLSource = "https://chapmanganato.com/manga-ng952689/chapter-700.5"

	err := mk.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 18, len(comic.Links))
}

func TestManganatoRetrieveIssueLinks(t *testing.T) {
	opt := &config.Options{
		URL:    "https://chapmanganato.com/manga-ng952689",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewManganato(opt)
	links, err := mk.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 748, len(links))
}
