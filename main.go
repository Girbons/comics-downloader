package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Girbons/comics-downloader/pkg/detector"
	"github.com/Girbons/comics-downloader/pkg/loader"
	log "github.com/sirupsen/logrus"
)

func usage() {
	fmt.Println("comics-download -url=http://comic-source")
	fmt.Println("comics-download -url=http://comic-source -debug")
	fmt.Println("comics-download -url=http://comic-source -country=italy")
	os.Exit(0)
}

func main() {
	// by default use log INFO Level
	log.SetLevel(log.InfoLevel)
	// setup the arguments
	url := flag.String("url", "", "Comic URL")
	debug := flag.Bool("debug", false, "Run the script in debug mode")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	// when you invoke `-- help` usage will appear
	flag.Usage = usage
	flag.Parse()

	// url is required
	if *url == "" {
		fmt.Println("url is required")
		os.Exit(1)
	}

	// if debug is true change log level to DEBUG
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// check if the url is supported
	source, check := detector.DetectComic(*url)
	if !check {
		log.Error("This site is not supported yet :(")
		os.Exit(1)
	}
	// in case the url is supported
	// setup the right strategy to parse a comic
	comic := loader.LoadComicFromSource(source, *url, *country)
	// make the PDF
	comic.MakeComic()
}
