package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComic(t *testing.T) {
	comic := new(Comic)
	// links
	links := []string{"foo.example.com"}
	// set info
	comic.SetInfo("foo", "2", "regex")
	// set links
	comic.SetImageLinks(links)
	// set the source
	comic.SetSource("bar")

	assert.Equal(t, "foo", comic.Name)
	assert.Equal(t, "2", comic.IssueNumber)
	assert.Equal(t, "regex", comic.ImageRegex)
	assert.Equal(t, "bar", comic.Source)

	assert.Equal(t, 1, len(comic.Links))
}
