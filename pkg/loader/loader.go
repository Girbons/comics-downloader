package loader

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites"
	log "github.com/sirupsen/logrus"
)

// LoadComicFromSource load the right comic strategy
func LoadComicFromSource(source, url string) *core.Comic {
	comic := new(core.Comic)

	comic.SetURLSource(url)
	comic.SetSource(source)

	switch source {
	case "www.comicextra.com":
		sites.SetupComicExtra(comic)
	case "www.mangahere.cc":
		sites.SetupMangaHere(comic)
	case "mangarock.com":
		sites.SetupMangaRock(comic)
	case "www.mangareader.net":
		sites.SetupMangaReader(comic)
	default:
		log.Warning("Cannot select a right strategy")
	}

	return comic
}
