package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/loader"
	log "github.com/sirupsen/logrus"
)

func main() {
	// use log INFO Level
	log.SetLevel(log.InfoLevel)
	// arguments setup
	url := flag.String("url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	format := flag.String("format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	flag.Parse()

	// url is required
	if *url == "" {
		fmt.Println("url parameter is required")
		os.Exit(1)
	}

	if !strings.HasSuffix(*url, ",") {
		*url = *url + ","
	}

	// check if the format is supported
	if !detector.DetectFormatOutput(*format) {
		log.WithFields(log.Fields{
			"format": *format,
		}).Error("Format not supported pdf will be used instead")
	}

	for _, u := range strings.Split(*url, ",") {
		if u != "" {
			// check if the url is supported
			source, check := detector.DetectComic(u)
			if !check {
				log.WithFields(log.Fields{"site": u}).Error("This site is not supported :(")
				continue
			}

			log.WithFields(log.Fields{"url": u}).Info("Downloading...")
			// in case the url is supported
			// setup the right strategy to parse a comic
			comic := loader.LoadComicFromSource(source, u, *country)
			comic.SetFormat(*format)
			comic.MakeComic()
		}
	}
}
