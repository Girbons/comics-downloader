package sites

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestMangarockSetup(t *testing.T) {
	comic := new(core.Comic)
	options := map[string]string{"country": "italy"}
	comic.SetOptions(options)
	comic.URLSource = "https://mangarock.com/manga/mrs-serie-35593/chapter/mrs-chapter-100051049"

	SetupMangaRock(comic)

	assert.Equal(t, "Boruto: Naruto Next Generations", comic.Name)
	assert.Equal(t, "Vol.4 Chapter 14: Teamwork...!!", comic.IssueNumber)
	assert.Equal(t, 49, len(comic.Links))
}
