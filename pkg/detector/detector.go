package detector

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// DetectComic will look for the url source in order to properly
// instantiate the Comic struct.
func DetectComic(url string) (string, bool) {
	var supportedSites = []string{"www.comicextra.com", "mangarock.com", "www.mangareader.net"}

	source, err := util.URLSource(url)

	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Error(err)
	}

	for _, site := range supportedSites {
		if source == site {
			return source, true
		}
	}

	return "", false
}

// DetectFormatOutput will check if the format is supported
func DetectFormatOutput(format string) bool {
	var supportedFormat = []string{core.PDF, core.EPUB, core.CBR, core.CBZ}
	return util.IsValueInSlice(format, supportedFormat)

}
