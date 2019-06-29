package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/sites"
	log "github.com/sirupsen/logrus"
)

var Messages = make(chan string)
var AppStatus = make(chan bool)

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

func GuiRun(link, format, country string, all, last bool) {
	// used for the gui to disable the `download` button
	if link == "" {
		msg := "url parameter required"
		log.Error(msg)
		sendToChannel(true, msg)
	}
	// used to disable submit button
	AppStatus <- true
	download(link, format, country, all, last, true)
	AppStatus <- false
}

// Run will run the downloader app
func Run(link, format, country string, all, last, deamon bool, sleepTime int) {
	if all && last {
		last = false
		log.Warning("all and last are selected, all parameter will be used")
	}

	// link is required
	if link == "" {
		msg := "url parameter required"
		log.Error(msg)
		os.Exit(1)
	}

	if deamon {
		for {
			fmt.Println("fooo")
			download(link, format, country, all, last, false)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			fmt.Println("bar")
		}
	}

	download(link, format, country, all, last, false)
}
