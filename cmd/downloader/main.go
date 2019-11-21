package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Girbons/comics-downloader/cmd/app"
	"github.com/Girbons/comics-downloader/internal/version"
)

func main() {
	all := flag.Bool("all", false, "Download all issues of the Comic or Comics")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	daemon := flag.Bool("daemon", false, "Run the download as daemon")
	format := flag.String("format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	imagesOnly := flag.Bool("images-only", false, "Download comic/manga images")
	imagesFormat := flag.String("images-format", "jpg", "To use with `images-only` flag, choose the image format, available png,jpeg,img")
	last := flag.Bool("last", false, "Download the last Comic issue")
	timeout := flag.Int("timeout", 600, "Timeout (seconds), specifies how often the downloader runs")
	url := flag.String("url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	versionFlag := flag.Bool("version", false, "Display release version")
	outputFolder := flag.String("output", "", "Folder where the comics will be saved")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("comics-downloader version %s", version.Tag)
		os.Exit(0)
	}

	// is this the best way?
	if *url == "" {
		for _, v := range flag.Args() {
			if !strings.HasPrefix(v, "-") || !strings.HasPrefix(v, "--") {
				if strings.HasPrefix(v, "http") || strings.HasPrefix(v, "https") {
					*url = *url + fmt.Sprintf("%s,", v)
				}
			}
		}
	}

	app.Run(*url, *format, *country, *imagesFormat, *all, *last, *daemon, *imagesOnly, *timeout, *outputFolder)
}
