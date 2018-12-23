package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func exists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func TestNewComic(t *testing.T) {
	comic := new(Comic)
	// links
	links := []string{"foo.example.com"}

	comic.SetInfo("foo", "2", "regex")
	comic.SetImageLinks(links)
	comic.SetSource("bar")

	assert.Equal(t, "foo", comic.Name)
	assert.Equal(t, "2", comic.IssueNumber)
	assert.Equal(t, "regex", comic.ImageRegex)
	assert.Equal(t, "bar", comic.Source)

	assert.Equal(t, 1, len(comic.Links))
}

func TestComicSetName(t *testing.T) {
	comic := new(Comic)

	comic.SetName("foo")

	assert.Equal(t, "foo", comic.Name)
}

func TestComicSetIssueNumber(t *testing.T) {
	comic := new(Comic)

	comic.SetIssueNumber("issue-number")

	assert.Equal(t, "issue-number", comic.IssueNumber)
}

func TestComicSetURLSource(t *testing.T) {
	comic := new(Comic)

	comic.SetURLSource("http://example.com")

	assert.Equal(t, "http://example.com", comic.URLSource)
}

func TestMakeComic(t *testing.T) {
	comic := new(Comic)

	comic.SetName("example.com")
	comic.SetIssueNumber("example-chapter-1")

	links := []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.SetImageLinks(links)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "example.com", "example-chapter-1.pdf"))

	assert.Equal(t, true, exists(dir))
}
