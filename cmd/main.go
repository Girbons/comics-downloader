package main

import (
	"flag"

	"github.com/Girbons/comics-downloader/cmd/app"
	log "github.com/sirupsen/logrus"
)

func init() {
	// use log INFO Level
	log.SetLevel(log.InfoLevel)
}

func main() {
	// arguments setup
	url := flag.String("url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	format := flag.String("format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	flag.Parse()

	app.Run(*url, *format, *country)
}
