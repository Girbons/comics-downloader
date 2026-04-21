package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"os"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/internal/version"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/detector"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
	"github.com/Girbons/comics-downloader/pkg/sites"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/sirupsen/logrus"
)

var (
	// AppStatus is used in GUI app to disable the `download` button
	AppStatus = make(chan bool)
	// Messages is used in GUI app to show app logs inside its specific box
	Messages = make(chan string)
)

// Runner orchestrates a download session based on immutable user input.
type Runner struct {
	base          config.Options
	bindToChannel bool
	messages      chan string

	loggerFactory func(bind bool, messages chan string) *logger.Logger
	clientFactory func() *httpclient.ComicClient
	sleep         func(time.Duration)
}

// NewRunner returns a Runner with default factories.
func NewRunner(base config.Options) *Runner {
	clientOpts := buildClientOptions(base)

	return &Runner{
		base: base,
		loggerFactory: func(bind bool, messages chan string) *logger.Logger {
			return logger.NewLogger(bind, messages)
		},
		clientFactory: func() *httpclient.ComicClient {
			return httpclient.NewComicClient(clientOpts...)
		},
		sleep: time.Sleep,
	}
}

// WithChannelBinding enables GUI-friendly logging via the provided channel.
func (r *Runner) WithChannelBinding(messages chan string) {
	r.bindToChannel = true
	r.messages = messages
}

func (r *Runner) prepareOptions() config.Options {
	opts := r.base
	if r.loggerFactory != nil {
		opts.Logger = r.loggerFactory(r.bindToChannel, r.messages)
	}
	if r.clientFactory != nil {
		opts.Client = r.clientFactory()
	}
	return opts
}

// Run executes a download session, respecting daemon configuration.
func (r *Runner) Run() {
	opts := r.prepareOptions()
	if opts.Logger == nil {
		opts.Logger = logger.NewLogger(false, nil)
	}
	if opts.Client == nil {
		opts.Client = httpclient.NewComicClient()
	}

	if opts.URL == "" {
		opts.Logger.Error("url parameter is required")
		return
	}

	if opts.Debug {
		opts.Logger.SetLevel(logrus.DebugLevel)
	}

	// daemon is started only if `all` or `last` flags are used
	if opts.Daemon && (opts.All || opts.Last) {
		for {
			r.download(opts)
			r.sleep(time.Duration(opts.DaemonTimeout) * time.Second)
		}
	}

	r.download(opts)
}

func (r *Runner) download(base config.Options) {
	opts := base

	if opts.All && opts.Last {
		opts.Last = false
		opts.Logger.Warning("all and last are selected, all parameter will be used")
	}

	// enforce `all` flag when `range` is used.
	if opts.IssuesRange != "" && !opts.All {
		opts.All = true
	}

	outputFolder := opts.OutputFolder
	if outputFolder == "" {
		dir, err := os.Getwd()
		if err != nil {
			opts.Logger.Errorf("Error determining current directory: %v", err)
			outputFolder = "."
		} else {
			outputFolder = dir
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	isNewVersionAvailable, newVersionLink, err := version.IsNewAvailable(ctx, opts.Client.HTTPClient())
	if err != nil {
		if errors.Is(err, version.ErrInvalidSemverTag) {
			opts.Logger.Debugf("Skipping version check: %v", err)
		} else {
			opts.Logger.Errorf("There was an error while checking for a new comics-downloader version: %v", err)
		}
	}

	if isNewVersionAvailable {
		opts.Logger.Infof("A new comics-downloader version is available at %s", newVersionLink)
	}

	rawURLs := strings.Split(opts.URL, ",")
	rawURLs = util.RemoveDuplicates(rawURLs)
	if len(rawURLs) > 1 {
		opts.Logger.Infof("Processing %d URLs", len(rawURLs))
	}

	for _, rawURL := range rawURLs {
		cleanedURL := strings.TrimSpace(rawURL)
		if cleanedURL == "" {
			continue
		}

		perURL := opts
		perURL.URL = cleanedURL
		perURL.OutputFolder = outputFolder
		// check if the link is supported
		source, isSupported, isDisabled := detector.DetectSource(cleanedURL)

		perURL.SourceName = source

		if !isSupported {
			perURL.Logger.Error("This site is not supported")
			continue
		}

		if isDisabled {
			perURL.Logger.Warning("Site currently disabled, please check https://github.com/Girbons/comics-downloader/issues/")
			continue
		}

		perURL.Logger.Info("Downloading... " + cleanedURL)
		collection, err := sites.LoadComicFromSource(&perURL)
		if err != nil {
			perURL.Logger.Error(err.Error())
			continue
		}

		for _, comic := range collection {
			if perURL.ImagesOnly {
				_, err = comic.DownloadImages(&perURL)
			} else {
				err = comic.MakeComic(&perURL)
			}

			if err != nil {
				perURL.Logger.Error(err.Error())
			}
		}
	}
}

// GuiRun will start the GUI app.
func GuiRun(options *config.Options) {
	AppStatus <- true
	runner := NewRunner(*options)
	runner.WithChannelBinding(Messages)
	runner.Run()
	AppStatus <- false
}

// Run will start the CLI app.
func Run(options *config.Options) {
	runner := NewRunner(*options)
	runner.Run()
}

func buildClientOptions(base config.Options) []httpclient.Option {
	defaultUA := fmt.Sprintf("comics-downloader/%s", version.Tag)
	agents := mergeUserAgents(defaultUA, base.UserAgents)
	opts := []httpclient.Option{httpclient.WithUserAgents(agents)}
	if !base.NoCache {
		opts = append(opts, httpclient.WithResponseCache(httpclient.NewInMemoryResponseCache(0)))
	}

	if strings.TrimSpace(base.SessionCookie) != "" {
		opts = append(opts, httpclient.WithHeaders(map[string]string{
			"Cookie": base.SessionCookie,
		}))
	}

	if proxyClient, ok := proxyHTTPClient(base.HTTPProxy); ok {
		opts = append(opts, httpclient.WithHTTPClient(proxyClient))
	}

	return opts
}

func proxyHTTPClient(rawProxy string) (*http.Client, bool) {
	trimmed := strings.TrimSpace(rawProxy)
	if trimmed == "" {
		return nil, false
	}

	proxyURL, err := urlpkg.Parse(trimmed)
	if err != nil || proxyURL.Scheme == "" || proxyURL.Host == "" {
		return nil, false
	}

	baseTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, false
	}

	transport := baseTransport.Clone()
	transport.Proxy = http.ProxyURL(proxyURL)

	return &http.Client{Transport: transport}, true
}

func mergeUserAgents(defaultAgent string, provided []string) []string {
	candidates := append([]string{defaultAgent}, provided...)
	seen := make(map[string]struct{}, len(candidates))
	var result []string
	for _, candidate := range candidates {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	if len(result) == 0 {
		result = append(result, defaultAgent)
	}
	return result
}
