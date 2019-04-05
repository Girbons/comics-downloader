package sites

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites/comicextra"
	"github.com/Girbons/comics-downloader/pkg/sites/mangareader"
	"github.com/Girbons/comics-downloader/pkg/sites/mangarock"
	"github.com/Girbons/comics-downloader/pkg/sites/mangatown"
)

func LoadComic(comic *core.Comic, country string) error {
	var err error

	switch comic.Source {
	case "www.comicextra.com":
		err = comicextra.Initialize(comic)
	case "mangarock.com":
		if country != "" {
			options := map[string]string{"country": country}
			comic.SetOptions(options)
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

	return err
}
