package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func usage() {
	fmt.Println("comics-download -url=http://comic-source")
	fmt.Println("comics-download -url=http://comic-source -debug")
	os.Exit(0)
}

func main() {
	log.SetLevel(log.InfoLevel)

	url := flag.String("url", "", "Comic URL")
	debug := flag.Bool("debug", false, "Run the script in debug mode")

	flag.Usage = usage
	flag.Parse()

	if *url == "" {
		fmt.Println("url parameter is required")
		os.Exit(1)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	DetectComic(*url)
}
