package app

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/sites"
	log "github.com/sirupsen/logrus"
)

func init() {
	// use log INFO Level
	log.SetLevel(log.InfoLevel)
}

// Run will run the downloader app
func Run(link string, format string, country string, all bool) {
	conf := new(config.ComicConfig)
	if err := conf.LoadConfig(); err != nil {
		log.Warning(err)
	}

	// link is required
	if link == "" {
		log.Error("url parameter is required")
	}

	// TODO: This doesn't seem necessary
	if !strings.HasSuffix(link, ",") {
		link = link + ","
	}

	for _, u := range strings.Split(link, ",") {
		if u != "" {
			// check if the link is supported
			source, check := detector.DetectComic(u)
			if !check {
				log.WithFields(log.Fields{"site": u}).Error("This site is not supported :(")
				continue
			}

			log.WithFields(log.Fields{
				"url": u,
			}).Info("Downloading...")
			// in case the link is supported
			// setup the right strategy to parse a comic
			collection, err := sites.LoadComicFromSource(conf, source, u, country, format, all)
			if err != nil {
				log.Error(err)
				continue
			}

			for _, comic := range collection {
				err = comic.MakeComic()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
