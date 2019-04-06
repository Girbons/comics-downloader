package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Girbons/mangarock"
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

	comic.SetInfo("foo", "2")
	comic.SetImageLinks(links)
	comic.SetSource("bar")

	assert.Equal(t, "foo", comic.Name)
	assert.Equal(t, "2", comic.IssueNumber)
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

func TestSetOptions(t *testing.T) {
	comic := new(Comic)

	options := map[string]string{"option": "foo"}
	comic.SetOptions(options)

	assert.Equal(t, comic.Options["option"], "foo")
}

func TestSplitURL(t *testing.T) {
	comic := new(Comic)
	comic.URLSource = "https://www.mangareader.net/naruto/1/"

	assert.Equal(t, comic.SplitURL()[3], "naruto")
	assert.Equal(t, comic.SplitURL()[4], "1")
}

func TestComicSetURLSource(t *testing.T) {
	comic := new(Comic)

	comic.SetURLSource("http://example.com")

	assert.Equal(t, "http://example.com", comic.URLSource)
}

func TestMakeComicPDF(t *testing.T) {
	comic := new(Comic)

	comic.SetName("example.com")
	comic.SetFormat("pdf")
	comic.SetIssueNumber("example-chapter-1")

	links := []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.SetImageLinks(links)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "example.com", "example-chapter-1.pdf"))

	assert.True(t, exists(dir))
}

func TestMakeComicEPUB(t *testing.T) {
	comic := new(Comic)

	comic.SetName("example.com")
	comic.SetFormat("epub")
	comic.SetIssueNumber("example-chapter-1")
	comic.SetAuthor("author")

	links := []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.SetImageLinks(links)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "example.com", "example-chapter-1.epub"))

	assert.True(t, exists(dir))
}

func TestMakeComicEPUBMangarock(t *testing.T) {
	client := mangarock.NewClient()
	options := map[string]string{"country": "italy"}
	client.SetOptions(options)
	result, _ := client.Pages("mrs-chapter-100051049")

	comic := new(Comic)

	comic.SetName("Boruto")
	comic.SetFormat("epub")
	comic.SetIssueNumber("chapter-13")
	comic.Source = "mangarock.com"
	comic.SetImageLinks(result.Data)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "mangarock.com", "Boruto", "chapter-13.epub"))

	assert.True(t, exists(dir))
}

func TestMakeComicCBZMangarock(t *testing.T) {
	client := mangarock.NewClient()
	options := map[string]string{"country": "italy"}
	client.SetOptions(options)
	result, _ := client.Pages("mrs-chapter-100051049")

	comic := new(Comic)

	comic.SetName("Boruto")
	comic.SetFormat("cbz")
	comic.SetIssueNumber("chapter-13")
	comic.Source = "mangarock.com"
	comic.SetImageLinks(result.Data)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "mangarock.com", "Boruto", "chapter-13.cbz"))

	assert.True(t, exists(dir))
}

func TestMakeComicCBRMangarock(t *testing.T) {
	client := mangarock.NewClient()
	options := map[string]string{"country": "italy"}
	client.SetOptions(options)
	result, _ := client.Pages("mrs-chapter-100051049")

	comic := new(Comic)

	comic.SetName("Boruto")
	comic.SetFormat("cbr")
	comic.SetIssueNumber("chapter-13")
	comic.Source = "mangarock.com"
	comic.SetImageLinks(result.Data)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "mangarock.com", "Boruto", "chapter-13.cbr"))

	assert.True(t, exists(dir))
}

func TestMakeComicPDFMangarock(t *testing.T) {
	client := mangarock.NewClient()
	options := map[string]string{"country": "italy"}
	client.SetOptions(options)
	result, _ := client.Pages("mrs-chapter-100051049")

	comic := new(Comic)

	comic.SetName("Boruto")
	comic.SetFormat("cbz")
	comic.SetIssueNumber("chapter-13")
	comic.Source = "mangarock.com"
	comic.SetImageLinks(result.Data)
	comic.MakeComic()

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "mangarock.com", "Boruto", "chapter-13.cbz"))

	assert.True(t, exists(dir))
}
