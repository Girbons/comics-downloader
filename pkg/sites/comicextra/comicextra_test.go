package comicextra

import (
	"testing"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestComicExtraSetup(t *testing.T) {
	comic := new(core.Comic)
	comic.URLSource = "https://www.comicextra.com/batman-2016/chapter-58/full"

	err := Initialize(comic)

	assert.Nil(t, err)
	assert.Equal(t, "batman-2016", comic.Name)
	assert.Equal(t, "chapter-58", comic.IssueNumber)
	assert.Equal(t, 24, len(comic.Links))
}

func TestRetrieveIssueLinks(t *testing.T) {
	issues, err := RetrieveIssueLinks("https://www.comicextra.com/comic/100-bullets")

	assert.Nil(t, err)
	assert.Equal(t, 100, len(issues))
}
