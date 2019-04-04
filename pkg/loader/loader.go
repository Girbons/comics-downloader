package loader

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites"
	log "github.com/sirupsen/logrus"
)

// LoadComicFromSource load the right comic strategy
func LoadComicFromSource(source, url, country string) (*core.Comic, error) {
	var err error

	comic := new(core.Comic)
	comic.SetURLSource(url)
	comic.SetSource(source)

	switch source {
	case "www.comicextra.com":
		err = sites.SetupComicExtra(comic)
	//case "www.mangahere.cc":
	//sites.SetupMangaHere(comic)
	case "mangarock.com":
		if country != "" {
			options := map[string]string{"country": country}
			comic.SetOptions(options)
		}
		err = sites.SetupMangaRock(comic)
	case "www.mangareader.net":
		err = sites.SetupMangaReader(comic)
	default:
		log.Error("Cannot select a right strategy")
		err = fmt.Errorf("It was not possible to determine the source")
	}

	return comic, err
}
