package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
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

	comic.Name = "foo"
	comic.IssueNumber = "2"
	comic.Links = links
	comic.Source = "bar"
	comic.ImagesFormat = "png"

	assert.Equal(t, "foo", comic.Name)
	assert.Equal(t, "2", comic.IssueNumber)
	assert.Equal(t, "bar", comic.Source)

	assert.Equal(t, 1, len(comic.Links))
}

func TestMakeComicPDF(t *testing.T) {
	comic := new(Comic)

	comic.Name = "foo"
	comic.Format = "pdf"
	comic.IssueNumber = "example-chapter-1"
	comic.Links = []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.ImagesFormat = "png"

	opt := &config.Options{
		OutputFolder:      filepath.Dir(os.Args[0]),
		CreateDefaultPath: true,
		Debug:             false,
		Logger:            logger.NewLogger(false, make(chan string)),
	}
	err := comic.MakeComic(opt)
	assert.Nil(t, err)

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "foo", "foo-example-chapter-1.pdf"))
	assert.True(t, exists(dir))
}

func TestMakeComicEPUB(t *testing.T) {
	comic := new(Comic)

	comic.Name = "foo"
	comic.Format = "epub"
	comic.IssueNumber = "example-chapter-1"
	comic.Author = "author"
	comic.ImagesFormat = "png"

	comic.Links = []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}

	opt := &config.Options{
		OutputFolder:      filepath.Dir(os.Args[0]),
		CreateDefaultPath: true,
		Debug:             false,
		Logger:            logger.NewLogger(false, make(chan string)),
	}

	err := comic.MakeComic(opt)
	assert.Nil(t, err)

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "foo", "foo-example-chapter-1.epub"))
	assert.True(t, exists(dir))
}

func TestDownloadImagesPNGFormat(t *testing.T) {
	comic := new(Comic)

	comic.Name = "foo-png"
	comic.Source = "fake"
	comic.IssueNumber = "example-chapter-1"
	comic.Links = []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.ImagesFormat = "png"

	opt := &config.Options{
		OutputFolder:      filepath.Dir(os.Args[0]),
		Debug:             false,
		CreateDefaultPath: true,
		Logger:            logger.NewLogger(false, make(chan string)),
	}
	dir, err := comic.DownloadImages(opt)
	files, _ := ioutil.ReadDir(dir)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))
}

func TestDownloadImagesJPGFormat(t *testing.T) {
	comic := new(Comic)

	comic.Name = "foo-jpg"
	comic.Source = "fake"
	comic.IssueNumber = "example-chapter-1"
	comic.Links = []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.ImagesFormat = "jpg"

	opt := &config.Options{
		OutputFolder:      filepath.Dir(os.Args[0]),
		CreateDefaultPath: true,
		Debug:             false,
		Logger:            logger.NewLogger(false, make(chan string)),
	}
	dir, err := comic.DownloadImages(opt)
	files, _ := ioutil.ReadDir(dir)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))
}

func TestDownloadImagesJPEGFormat(t *testing.T) {
	comic := new(Comic)

	comic.Name = "bar-jpeg"
	comic.Source = "fake"
	comic.IssueNumber = "example-chapter-1"
	comic.ImagesFormat = "jpeg"
	comic.Links = []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}

	opt := &config.Options{
		OutputFolder:      filepath.Dir(os.Args[0]),
		CreateDefaultPath: true,
		Debug:             false,
		Logger:            logger.NewLogger(false, make(chan string)),
	}
	dir, err := comic.DownloadImages(opt)
	files, _ := ioutil.ReadDir(dir)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))
}

func TestDownloadImagesIMGFormat(t *testing.T) {
	comic := new(Comic)

	comic.Name = "bar-img"
	comic.Source = "fake"
	comic.IssueNumber = "example-chapter-1"
	comic.Links = []string{"http://via.placeholder.com/150", "http://via.placeholder.com/150", "http://via.placeholder.com/150"}
	comic.ImagesFormat = "img"

	opt := &config.Options{
		OutputFolder:      filepath.Dir(os.Args[0]),
		CreateDefaultPath: true,
		Debug:             false,
		Logger:            logger.NewLogger(false, make(chan string)),
	}
	dir, err := comic.DownloadImages(opt)
	files, _ := ioutil.ReadDir(dir)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))
}
