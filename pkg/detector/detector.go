package detector

import (
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// SupportedSites are the supported sites.
var SupportedSites = map[string]map[string]bool{
	"www.comicextra.com":  {"isDisabled": false},
	"readcomiconline.to":  {"isDisabled": false},
	"www.mangareader.net": {"isDisabled": false},
	"www.mangatown.com":   {"isDisabled": false},
	"www.mangahere.cc":    {"isDisabled": false},
	"mangadex.cc":         {"isDisabled": false},
	"mangadex.org":        {"isDisabled": false},
}

// DetectComic will look for the url source to check if a source is supported.
func DetectComic(url string) (string, bool, bool) {
	var (
		isSupported bool
		isDisabled  bool
		source      string
	)

	source, err := util.URLSource(url)

	if err != nil {
		log.Error(err)
	}

	for k, v := range SupportedSites {
		if k != source {
			continue
		}

		isSupported = true
		isDisabled = v["isDisabled"]
	}

	return source, isSupported, isDisabled
}
