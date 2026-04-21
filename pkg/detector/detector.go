package detector

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/sites"
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

// DetectSource will look for the url source to check if a source is supported.
func DetectSource(url string) (source string, isSupported, isDisabled bool) {
	source, err := util.URLSource(url)

	if err != nil {
		log.Error(err)
	}

	for siteName, supportedSite := range sites.SupportedSites {
		if !strings.Contains(source, siteName) {
			continue
		}

		isSupported = true
		isDisabled = !supportedSite.IsEnabled
	}

	return source, isSupported, isDisabled
}
