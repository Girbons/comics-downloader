package app

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/sites"
	log "github.com/sirupsen/logrus"
)

var Messages = make(chan string)

func init() {
	// use log INFO Level
	log.SetLevel(log.InfoLevel)
}

func sendToChannel(enabled bool, message string) {
	if enabled {
		Messages <- message
	}
}

// Run will run the downloader app
func Run(link, format, country string, all, bindLogsToChannel bool) {
	conf := new(config.ComicConfig)
	if err := conf.LoadConfig(); err != nil {
		log.Warning(err)
	}

	// link is required
	if link == "" {
		msg := "url parameter required"
		log.Error(msg)
		sendToChannel(bindLogsToChannel, msg)
	}

	for _, u := range strings.Split(link, ",") {
		if u != "" {
			// check if the link is supported
			source, check := detector.DetectComic(u)
			if !check {
				msg := "This site is not supported :("
				log.WithFields(log.Fields{"site": u}).Error(msg)
				sendToChannel(bindLogsToChannel, msg)
				continue
			}

			msg := "Downloading..."
			log.WithFields(log.Fields{"url": u}).Info(msg)
			sendToChannel(bindLogsToChannel, msg)
			// in case the link is supported
			// setup the right strategy to parse a comic
			collection, err := sites.LoadComicFromSource(conf, source, u, country, format, all)
			if err != nil {
				log.Error(err)
				sendToChannel(bindLogsToChannel, fmt.Sprintf("ERROR: %s", err))
				continue
			}

			for _, comic := range collection {
				err = comic.MakeComic()
				if err != nil {
					log.Error(err)
					sendToChannel(bindLogsToChannel, fmt.Sprintf("ERROR: %s", err))
				} else {
					sendToChannel(bindLogsToChannel, fmt.Sprintf("%s, Succesfully Downloaded", comic.Name))
				}
			}
		}
	}
}
