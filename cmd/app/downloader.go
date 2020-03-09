package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/config"
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

func download(options *config.Options, bindLogsToChannel bool) {
	if options.OutputFolder == "" {
		options.OutputFolder = filepath.Dir(os.Args[0])
	}

	for _, u := range strings.Split(options.Url, ",") {
		if u != "" {
			// check if the link is supported
			source, check, isDisabled := detector.DetectComic(u)

			options.Source = source

			if !check {
				msg := "This site is not supported :("
				log.WithFields(log.Fields{"site": u}).Error(msg)
				sendToChannel(bindLogsToChannel, msg)
				continue
			}

			if isDisabled {
				msg := "Site currently disabled, please check https://github.com/Girbons/comics-downloader/issues/"
				log.WithFields(log.Fields{"site": u}).Warning(msg)
				sendToChannel(bindLogsToChannel, msg)
				continue
			}

			msg := "Downloading..."
			log.WithFields(log.Fields{"url": u}).Info(msg)
			sendToChannel(bindLogsToChannel, msg)
			// in case the link is supported
			// setup the right strategy to parse a comic
			collection, err := sites.LoadComicFromSource(options)
			if err != nil {
				log.Error(err)
				sendToChannel(bindLogsToChannel, fmt.Sprintf("ERROR: %s", err))
				continue
			}

			for _, comic := range collection {
				if options.ImagesOnly {
					_, err = comic.DownloadImages(options.OutputFolder)
				} else {
					err = comic.MakeComic(options.OutputFolder)
				}
				checkErr(err, bindLogsToChannel, comic)
			}
		}
	}
}

// GuiRun will start the GUI app
func GuiRun(options *config.Options) {
	AppStatus <- true
	download(options, true)
	AppStatus <- false
}

// Run will start the CLI app
func Run(options *config.Options) {
	if options.All && options.Last {
		options.Last = false
		log.Warning("all and last are selected, all parameter will be used")
	}

	// link is required
	if options.Url == "" {
		log.Fatal("url parameter is required")
	}

	// daemon is started only if `all` or `last` flags are used
	if options.Daemon && (options.All || options.Last) {
		for {
			download(options, false)
			time.Sleep(time.Duration(options.Timeout) * time.Second)
		}
	}

	download(options, false)
}
