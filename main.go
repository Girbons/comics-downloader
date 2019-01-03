package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/loader"
	log "github.com/sirupsen/logrus"
)

func main() {
	// use log INFO Level
	log.SetLevel(log.InfoLevel)
	// arguments setup
	url := flag.String("url", "", "Comic URL")
	format := flag.String("format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	flag.Parse()

	// url is required
	if *url == "" {
		fmt.Println("url parameter is required")
		os.Exit(1)
	}

	// check if the url is supported
	source, check := detector.DetectComic(*url)
	if !check {
		log.WithFields(log.Fields{
			"site": *url,
		}).Error("This site is not supported :(")
		os.Exit(1)
	}

	// check if the format is supported
	if !detector.DetectFormatOutput(*format) {
		log.WithFields(log.Fields{
			"format": *format,
		}).Info("Format not supported pdf will be used instead")
	}

	// in case the url is supported
	// setup the right strategy to parse a comic
	comic := loader.LoadComicFromSource(source, *url, *country)
	comic.SetFormat(*format)
	comic.MakeComic()
}
