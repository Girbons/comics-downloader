package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Girbons/comics-downloader/cmd/app"
	"github.com/Girbons/comics-downloader/internal/version"
	"github.com/Girbons/comics-downloader/pkg/config"
)

var (
	// shows debug log
	debug bool
	// download mode options
	all  bool
	last bool
	// images options
	imagesOnly   bool
	imagesFormat string
	// manga/comic language
	country string
	// manga/comic final output
	format string
	// url source
	url string
	// folder where the files will be saved
	outputFolder string
	// daemon options
	daemon  bool
	timeout int
	// app version
	versionFlag bool
	// range of issues to download
	issuesRange string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Shows Debug log")
	flag.BoolVar(&all, "all", false, "Download all issues of the Comic or Comics")
	flag.BoolVar(&daemon, "daemon", false, "Run the download as daemon")
	flag.BoolVar(&imagesOnly, "images-only", false, "Download comic/manga images")
	flag.BoolVar(&last, "last", false, "Download the last Comic issue")
	flag.BoolVar(&versionFlag, "version", false, "Display release version")

	flag.StringVar(&country, "country", "", "Set the country to retrieve a manga, Used by MangaDex which uses ISO 3166-1 codes")
	flag.StringVar(&format, "format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	flag.StringVar(&imagesFormat, "images-format", "jpg", "To use with `images-only` flag, choose the image format, available png,jpeg,img")
	flag.StringVar(&url, "url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	flag.StringVar(&outputFolder, "output", "", "Folder where the comics will be saved")
	flag.StringVar(&issuesRange, "range", "", "Range of issues to download, example 3-9")

	flag.IntVar(&timeout, "timeout", 600, "Timeout (seconds), specifies how often the downloader runs")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Println("comics-downloader version", version.Tag)
		os.Exit(0)
	}

	options := &config.Options{
		Debug:        debug,
		All:          all,
		Last:         last,
		Country:      country,
		ImagesOnly:   imagesOnly,
		ImagesFormat: imagesFormat,
		URL:          url,
		Format:       format,
		Daemon:       daemon,
		Timeout:      timeout,
		OutputFolder: outputFolder,
		IssuesRange:  issuesRange,
	}

	app.Run(options)
}
