package detector

import (
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// SupportedSites are the supported sites.
var SupportedSites = map[string]map[string]bool{
	"mangadex.org":       {"isDisabled": true},
	"mangareader.tv":     {"isDisabled": true},
	"readallcomics.com":  {"isDisabled": false},
	"readcomiconline.li": {"isDisabled": false},
	"www.comicextra.com": {"isDisabled": false},
	"www.mangahere.cc":   {"isDisabled": false},
	"www.mangatown.com":  {"isDisabled": false},
	"mangakakalot.com":   {"isDisabled": false},
	"manganato.com":      {"isDisabled": false},
	"readmanganato.com":  {"isDisabled": false},
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
