package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/internal/version"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/http"
	"github.com/Girbons/comics-downloader/pkg/sites"
	"github.com/sirupsen/logrus"
)

var (
	// AppStatus is used in GUI app to disable the `download` button
	AppStatus = make(chan bool)
	// Messages is used in GUI app to show app logs inside its specific box
	Messages = make(chan string)
)

func download(options *config.Options) {
	if options.Debug {
		options.Logger.SetLevel(logrus.DebugLevel)
	}

	if options.All && options.Last {
		options.Last = false
		options.Logger.Warning("all and last are selected, all parameter will be used")
	}

	// enforce `all` flag when `range` is used.
	if options.IssuesRange != "" && !options.All {
		options.All = true
	}

	if options.OutputFolder == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining current directory: %v\n", err)
			options.OutputFolder = "."
		} else {
			options.OutputFolder = dir
		}
	}

	isNewVersionAvailable, newVersionLink, err := version.IsNewAvailable()
	if err != nil {
		options.Logger.Error("There was an error while checking for a new comics-downloader version")
	}

	if isNewVersionAvailable {
		options.Logger.Info(fmt.Sprintf("A new comics-downloader version is available at %s", newVersionLink))
	}

	urls := options.URL

	for _, u := range strings.Split(urls, ",") {
		if u == "" {
			continue
		}

		// check if the link is supported
		source, check, isDisabled := detector.DetectComic(u)

		options.Source = source
		options.URL = u

		if !check {
			options.Logger.Error("This site is not supported")
			continue
		}

		if isDisabled {
			options.Logger.Warning("Site currently disabled, please check https://github.com/Girbons/comics-downloader/issues/")
			continue
		}

		options.Logger.Info("Downloading...")
		collection, err := sites.LoadComicFromSource(options)
		if err != nil {
			options.Logger.Error(err.Error())
			continue
		}

		for _, comic := range collection {
			if options.ImagesOnly {
				_, err = comic.DownloadImages(options)
			} else {
				err = comic.MakeComic(options)
			}

			if err != nil {
				options.Logger.Error(err.Error())
			}
		}
	}
}

// GuiRun will start the GUI app
func GuiRun(options *config.Options) {
	AppStatus <- true
	options.Logger = logger.NewLogger(true, Messages)
	download(options)
	AppStatus <- false
}

// Run will start the CLI app
func Run(options *config.Options) {
	options.Logger = logger.NewLogger(false, Messages)
	options.Client = http.NewComicClient()

	// link is required
	if options.URL == "" {
		options.Logger.Error("url parameter is required")
		return
	}

	// daemon is started only if `all` or `last` flags are used
	if options.Daemon && (options.All || options.Last) {
		for {
			download(options)
			time.Sleep(time.Duration(options.DaemonTimeout) * time.Second)
		}
	}

	download(options)
}
