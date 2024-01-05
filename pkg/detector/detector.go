package detector

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// SupportedSites are the supported sites.
var SupportedSites = map[string]map[string]bool{
	"comicextra":      {"isDisabled": false},
	"mangadex":        {"isDisabled": false},
	"mangareader":     {"isDisabled": true},
	"mangakakalot":    {"isDisabled": false},
	"manganato":       {"isDisabled": false},
	"mangatown":       {"isDisabled": false},
	"readallcomics":   {"isDisabled": false},
	"readcomiconline": {"isDisabled": false},
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
		if !strings.Contains(source, k) {
			continue
		}

		isSupported = true
		isDisabled = v["isDisabled"]
	}

	return source, isSupported, isDisabled
}
