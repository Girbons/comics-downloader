package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/sites"
	log "github.com/sirupsen/logrus"
)

var (
	// AppStatus is used in GUI app to disable the `download` button
	AppStatus = make(chan bool)
	// Messages is used in GUI app to show app logs inside its specific box
	Messages = make(chan string)
)

func init() {
	// use log INFO Level
	log.SetLevel(log.InfoLevel)
}

func sendToChannel(enabled bool, message string) {
	if enabled {
		Messages <- message
	}
}

func download(link, format, country string, all, last, bindLogsToChannel bool) {
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
			collection, err := sites.LoadComicFromSource(source, u, country, format, all, last)
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

// GuiRun will start the GUI app
func GuiRun(link, format, country string, all, last bool) {
	AppStatus <- true
	download(link, format, country, all, last, true)
	AppStatus <- false
}

// Run will start the CLI app
func Run(link, format, country string, all, last, deamon bool, timeout int) {
	if all && last {
		last = false
		log.Warning("all and last are selected, all parameter will be used")
	}

	// link is required
	if link == "" {
		log.Fatal("url parameter is required")
	}

	// deamon is started only if `all` or `last` flags are used
	if deamon && (all || last) {
		for {
			download(link, format, country, all, last, false)
			time.Sleep(time.Duration(timeout) * time.Second)
		}
	} else {
		log.Warning("To use `-deamon` be sure to pass `-all` or `-last` flags")
	}

	download(link, format, country, all, last, false)
}
