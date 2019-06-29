package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Girbons/comics-downloader/cmd/app"
)

var (
	release   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	all := flag.Bool("all", false, "Download all issues of the Comic or Comics")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	deamon := flag.Bool("deamon", false, "Run the download as deamon")
	format := flag.String("format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	last := flag.Bool("last", false, "Download the last Comic issue")
	sleepTime := flag.Int("sleepTime", 30, "Deamon sleep time")
	url := flag.String("url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	versionFlag := flag.Bool("version", false, "Display the release")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Built on %s from release %s\n", buildTime, release)
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

	app.Run(*url, *format, *country, *all, *last, *deamon, *sleepTime)
}
