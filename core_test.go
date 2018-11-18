package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUrlValid(t *testing.T) {
	validUrl := IsUrlValid("http://example.com")
	gifUrl := IsUrlValid("http://foo.gif")
	logoUrl := IsUrlValid("http://foo.logo")

	assert.Equal(t, true, validUrl)
	assert.Equal(t, false, gifUrl)
	assert.Equal(t, false, logoUrl)
}

func TestNewComic(t *testing.T) {
	comic := NewComic("foo", "2", "regex", "source")
	assert.Equal(t, "foo", comic.Name)
	assert.Equal(t, "2", comic.IssueNumber)
	assert.Equal(t, "regex", comic.ImageRegex)
	assert.Equal(t, "source", comic.Source)
}
