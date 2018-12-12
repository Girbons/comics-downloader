package detector

import (
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// DetectComic will look for the url source in order to properly
// instantiate the Comic struct.
func DetectComic(url string) (string, bool) {
	var supportedSites = []string{"www.comicextra.com", "www.mangahere.cc", "mangarock.com", "www.mangareader.net"}

	log.Debug("Detecting the source...")
	source, err := util.UrlSource(url)

	if err != nil {
		log.Error(err)
	}

	for _, site := range supportedSites {
		if source == site {
			log.Debug("Source detected")
			return source, true
		}
	}
	log.Debug("Source not detected: ", source)
	return "", false
}
