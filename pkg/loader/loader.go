package loader

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites"
)

// LoadComicFromSource create a comic instance
func LoadComicFromSource(source, url, country string) (*core.Comic, error) {
	comic := new(core.Comic)
	comic.SetURLSource(url)
	comic.SetSource(source)
	err := sites.LoadComic(comic, country)

	return comic, err
}
