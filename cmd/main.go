package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/cmd/app"
)

func main() {
	url := flag.String("url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	format := flag.String("format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	country := flag.String("country", "", "Set the country to retrieve a manga, Used by MangaRock")
	flag.Parse()

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

	app.Run(*url, *format, *country)
}
