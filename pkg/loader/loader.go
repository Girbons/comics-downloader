package loader

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites"
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// LoadComicFromSource load the right comic strategy
func LoadComicFromSource(source, url string) *core.Comic {
	comic := new(core.Comic)
	comic.SetURLSource(url)
	comic.SetSource(source)
	// split the url
	splittedUrl := util.SplitUrl(url)

	switch source {
	case "www.comicextra.com":
		sites.SetupComicExtra(comic, splittedUrl)
	case "www.mangahere.cc":
		sites.SetupMangaHere(comic, splittedUrl)
	case "mangarock.com":
		sites.SetupMangaRock(comic, splittedUrl)
	default:
		log.Warning("Cannot select a right strategy")
	}

	return comic
}
