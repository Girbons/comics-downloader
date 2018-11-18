package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TestDetectComicExtra(t *testing.T) {
	DetectComic("http://www.comicextra.com/daredevil-2016/chapter-600/full")

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "www.comicextra.com", "daredevil-2016"))
	result, _ := exists(dir)

	assert.Equal(t, true, result)
}

func TestDetectMangaHere(t *testing.T) {
	DetectComic("http://www.mangahere.cc/manga/shingeki_no_kyojin_before_the_fall/c048/")

	dir, _ := filepath.Abs(fmt.Sprintf("%s/%s/%s/%s/", filepath.Dir(os.Args[0]), "comics", "www.mangahere.cc", "shingeki_no_kyojin_before_the_fall"))
	result, _ := exists(dir)

	assert.Equal(t, true, result)
}
