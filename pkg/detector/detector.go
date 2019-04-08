package detector

import (
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// DetectComic will look for the url source in order to properly
// instantiate the Comic struct.
func DetectComic(url string) (string, bool) {
	var supportedSites = []string{
		"www.comicextra.com",
		"mangarock.com",
		"www.mangareader.net",
		"www.mangatown.com",
		"www.mangahere.cc",
	}

	source, err := util.URLSource(url)

	if err != nil {
		log.Error(err)
	}

	for _, site := range supportedSites {
		if source == site {
			return source, true
		}
	}

	return "", false
}
