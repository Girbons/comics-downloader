package sites

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites/comicextra"
	"github.com/Girbons/comics-downloader/pkg/sites/mangareader"
	"github.com/Girbons/comics-downloader/pkg/sites/mangarock"
	"github.com/Girbons/comics-downloader/pkg/sites/mangatown"
)

// LoadComicFromSource will return a `comic` instance initialized based on the source
func LoadComicFromSource(conf *config.ComicConfig, source, url, country, format string) (*core.Comic, error) {
	var err error

	comic := &core.Comic{
		Config:    conf,
		URLSource: url,
		Source:    source,
		Format:    format,
	}

	switch source {
	case "www.comicextra.com":
		err = comicextra.Initialize(comic)
	case "mangarock.com":
		if country != "" {
			options := map[string]string{"country": country}
			comic.Options = options
		}
		err = mangarock.Initialize(comic)
	case "www.mangareader.net":
		err = mangareader.Initialize(comic)
	case "www.mangatown.com":
		err = mangatown.Initialize(comic)
	case "www.mangahere.cc":
		err = fmt.Errorf("mangahere is currently disabled")
	//sites.SetupMangaHere(comic)
	default:
		err = fmt.Errorf("It was not possible to determine the source")
	}

	return comic, err
}
