package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

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
	forceAspect     bool
	format          string
	customComicName string
	// force only issue number filenames
	issueNumberNameOnly bool
	// url source
	url string
	// folder where the files will be saved
	outputFolder string
	// enable/disable default folder structure
	createDefaultPath bool
	// daemon options
	daemon        bool
	daemonTimeout int
	// app version
	versionFlag bool
	// range of issues to download
	issuesRange string
	// string to be used for each issue/chapter folder
	issueFolderName string
	// request customization
	userAgentsCSV string
	sessionCookie string
	// throttling
	requestDelay       time.Duration
	requestDelayJitter time.Duration
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
	flag.StringVar(&customComicName, "custom-comic-name", "", "Use a custom name for the comic output.")
	flag.StringVar(&imagesFormat, "images-format", "jpg", "To use with `images-only` flag, choose the image format, available png,jpeg,img")
	flag.BoolVar(&issueNumberNameOnly, "issue-number-only", false, "Force only saving with issue number instead of chapter name + issue number.")
	flag.StringVar(&url, "url", "", "Comic URL or Comic URLS by separating each site with a comma without the use of spaces")
	flag.StringVar(&outputFolder, "output", "", "Folder where the comics will be saved")
	flag.StringVar(&issuesRange, "range", "", "Range of issues to download, example 3-9")
	flag.StringVar(&issueFolderName, "issue-folder-name", "issue-", "Folder name where each issue/chapter will be saved, default 'issue-#'")
	flag.StringVar(&userAgentsCSV, "user-agents", "", "Comma-separated list of alternative User-Agent values to rotate per request")
	flag.StringVar(&sessionCookie, "session-cookie", "", "Custom Cookie header value (e.g., cf_clearance=...; other=...) for protected sources")
	flag.DurationVar(&requestDelay, "request-delay", config.DefaultRequestDelay, "Base delay inserted before downloading each image (e.g. 500ms)")
	flag.DurationVar(&requestDelayJitter, "request-delay-jitter", config.DefaultRequestDelayJitter, "Maximum additional random delay added to the base request delay (e.g. 250ms)")

	flag.IntVar(&daemonTimeout, "daemon-timeout", 600, "DaemonTimeout (seconds), specifies how often the downloader runs")
}

func buildOptions() config.Options {
	return config.Options{
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
		CustomComicName:     customComicName,
		Daemon:              daemon,
		DaemonTimeout:       daemonTimeout,
		OutputFolder:        outputFolder,
		CreateDefaultPath:   createDefaultPath,
		IssuesRange:         issuesRange,
		IssueFolderName:     issueFolderName,
		UserAgents:          splitAndTrim(userAgentsCSV),
		SessionCookie:       strings.TrimSpace(sessionCookie),
		RequestDelay:        requestDelay,
		RequestDelayJitter:  requestDelayJitter,
	}
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Println("comics-downloader version", version.Tag)
		os.Exit(0)
	}

	opts := buildOptions()
	app.Run(&opts)
}

func splitAndTrim(csv string) []string {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	var out []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
