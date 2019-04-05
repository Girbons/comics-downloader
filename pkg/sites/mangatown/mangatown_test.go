package mangatown

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestComicExtraSetup(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = "http://www.mangatown.com/manga/naruto/v63/c684/"

	err := Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, "naruto", comic.Name)
	assert.Equal(t, "c684", comic.IssueNumber)
	assert.Equal(t, 22, len(comic.Links))
}
