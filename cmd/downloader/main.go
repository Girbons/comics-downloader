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
	forceAspect bool
	format      string
	// force only issue number filenames
	issueNumberNameOnly bool
	// url source
	url string
	// folder where the files will be saved
	outputFolder string
	// enable/disable default folder structure
	createDefaultPath bool
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
	flag.BoolVar(&createDefaultPath, "create-default-path", true, "Using this flag your comics/issue will be downloaded without prepending the default folder structure, `comics/[source]/[name]/`")
	flag.StringVar(&country, "country", "", "Set the country to retrieve a manga, Used by MangaDex which uses ISO 3166-1 codes")
	flag.BoolVar(&forceAspect, "force-aspect", false, "Force images to A4 Portrait aspect ratio")
	flag.StringVar(&format, "format", "pdf", "Comic format output, supported formats are pdf,epub,cbr,cbz")
	flag.StringVar(&imagesFormat, "images-format", "jpg", "To use with `images-only` flag, choose the image format, available png,jpeg,img")
	flag.BoolVar(&issueNumberNameOnly, "issue-number-only", false, "Force only saving with issue number instead of chapter name + issue number.")
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
		Debug:               debug,
		All:                 all,
		Last:                last,
		Country:             country,
		ImagesOnly:          imagesOnly,
		ImagesFormat:        imagesFormat,
		IssueNumberNameOnly: issueNumberNameOnly,
		URL:                 url,
		ForceAspect:         forceAspect,
		Format:              format,
		Daemon:              daemon,
		Timeout:             timeout,
		OutputFolder:        outputFolder,
		CreateDefaultPath:   createDefaultPath,
		IssuesRange:         issuesRange,
	}

	app.Run(options)
}
