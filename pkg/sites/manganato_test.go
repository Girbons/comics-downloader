package sites

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestManganatoGetInfo(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readmanganato.com/manga-ny955333",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mg := NewManganato(opt)
	// readmanganato.com
	name, issueNumber := mg.GetInfo("https://readmanganato.com/manga-ny955333/chapter-36")
	assert.Equal(t, "To The Past", name)
	assert.Equal(t, "36", issueNumber)

	// manganato.com
	opt.URL = "https://manganato.com/manga-gd983838"
	name, issueNumber = mg.GetInfo("https://readmanganato.com/manga-gd983838/chapter-76")
	assert.Equal(t, "Voice And Written Words", name)
	assert.Equal(t, "76", issueNumber)
}

func TestManganatoSetup(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readmanganato.com/manga-ny955333",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewManganato(opt)
	comic := new(core.Comic)
	comic.URLSource = "https://readmanganato.com/manga-ny955333/chapter-36"

	err := mk.Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, 20, len(comic.Links))
}

func TestManganatoRetrieveIssueLinks(t *testing.T) {
	opt := &config.Options{
		URL:    "https://readmanganato.com/manga-ny955333",
		All:    false,
		Last:   false,
		Debug:  false,
		Logger: logger.NewLogger(false, make(chan string)),
	}
	mk := NewManganato(opt)
	links, err := mk.RetrieveIssueLinks()

	assert.Nil(t, err)
	assert.Equal(t, 58, len(links))
}
