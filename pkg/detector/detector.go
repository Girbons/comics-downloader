package detector

import (
	"github.com/Girbons/comics-downloader/pkg/sites"
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// DetectComic will look for the url source to check if a source is supported.
func DetectComic(url string) (string, bool) {
	source, err := util.URLSource(url)

	if err != nil {
		log.Error(err)
	}

	for _, site := range sites.SupportedSites {
		if source == site {
			return source, true
		}
	}

	return "", false
}
