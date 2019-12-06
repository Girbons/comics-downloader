package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/core"
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

func checkErr(err error, bindLogsToChannel bool, comic *core.Comic) {
	if err != nil {
		log.Error(err)
		sendToChannel(bindLogsToChannel, fmt.Sprintf("ERROR: %s", err))
	} else {
		name := fmt.Sprintf("%s %s.%s", comic.Name, comic.IssueNumber, comic.Format)
		sendToChannel(bindLogsToChannel, fmt.Sprintf("%s, Succesfully Downloaded", name))
	}
}

func download(link, format, country string, all, last, bindLogsToChannel, imagesOnly bool, imagesFormat, outputFolder string) {
	if outputFolder == "" {
		outputFolder = filepath.Dir(os.Args[0])
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

			if strings.Contains(source, "mangadex") && (last) {
				msg := "`last` parameters are not supported"
				log.WithFields(log.Fields{"site": u}).Warning(msg)
				sendToChannel(bindLogsToChannel, msg)
				last = false
			}

			msg := "Downloading..."
			log.WithFields(log.Fields{"url": u}).Info(msg)
			sendToChannel(bindLogsToChannel, msg)
			// in case the link is supported
			// setup the right strategy to parse a comic
			collection, err := sites.LoadComicFromSource(source, u, country, format, imagesFormat, all, last, imagesOnly, outputFolder)
			if err != nil {
				log.Error(err)
				sendToChannel(bindLogsToChannel, fmt.Sprintf("ERROR: %s", err))
				continue
			}

			for _, comic := range collection {
				if imagesOnly {
					_, err = comic.DownloadImages(outputFolder)
				} else {
					err = comic.MakeComic(outputFolder)
				}
				checkErr(err, bindLogsToChannel, comic)
			}
		}
	}
}

// GuiRun will start the GUI app
func GuiRun(link, format, country, imagesFormat string, all, last, imagesOnly bool, outputFolder string) {
	AppStatus <- true
	download(link, format, country, all, last, true, imagesOnly, imagesFormat, outputFolder)
	AppStatus <- false
}

// Run will start the CLI app
func Run(link, format, country, imagesFormat string, all, last, daemon, imagesOnly bool, timeout int, outputFolder string) {
	if all && last {
		last = false
		log.Warning("all and last are selected, all parameter will be used")
	}

	// link is required
	if link == "" {
		log.Fatal("url parameter is required")
	}

	// daemon is started only if `all` or `last` flags are used
	if daemon && (all || last) {
		for {
			download(link, format, country, all, last, false, imagesOnly, imagesFormat, outputFolder)
			time.Sleep(time.Duration(timeout) * time.Second)
		}
	}

	download(link, format, country, all, last, false, imagesOnly, imagesFormat, outputFolder)
}
