package app

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/loader"
	log "github.com/sirupsen/logrus"
)

func Run(link, format, country string) {
	// link is required
	if link == "" {
		log.Fatal("url parameter is required")
	}

	if !strings.HasSuffix(link, ",") {
		link = link + ","
	}

	// check if the format is supported
	if !detector.DetectFormatOutput(format) {
		log.WithFields(log.Fields{
			"format": format,
		}).Error("Format not supported PDF will be used instead")
	}

	for _, u := range strings.Split(link, ",") {
		if u != "" {
			// check if the link is supported
			source, check := detector.DetectComic(u)
			if !check {
				log.WithFields(log.Fields{"site": u}).Error("This site is not supported :(")
				continue
			}

			log.WithFields(log.Fields{"link": u}).Info("Downloading...")
			// in case the link is supported
			// setup the right strategy to parse a comic
			comic, err := loader.LoadComicFromSource(source, u, country)
			if err != nil {
				log.WithFields(log.Fields{"link": u}).Error(err)
				continue
			}
			comic.SetFormat(format)
			comic.MakeComic()
		}
	}
}
