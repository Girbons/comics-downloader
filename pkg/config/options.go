package config

import (
	"time"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/http"
)

const (
	// DefaultRequestDelay defines the base time to wait between subsequent image requests.
	DefaultRequestDelay = 500 * time.Millisecond
	// DefaultRequestDelayJitter adds up to this much random extra delay to avoid fixed patterns.
	DefaultRequestDelayJitter = 250 * time.Millisecond
	// DefaultRequestTimeout is the default timeout for HTTP requests.
	DefaulltRequestTimeout = 30 * time.Second
)

// Options represents the comics downloader options.
type Options struct {
	Debug               bool
	All                 bool
	Last                bool
	ImagesOnly          bool
	Daemon              bool
	DaemonTimeout       int
	ImagesFormat        string
	Country             string
	OutputFormat        string
	CustomComicName     string
	ForceAspect         bool
	OutputFolder        string
	CreateDefaultPath   bool
	IssueNumberNameOnly bool
	URL                 string
	SourceName          string
	IssuesRange         string
	IssueFolderName     string

	UserAgents         []string
	SessionCookie      string
	NoCache            bool
	RequestDelay       time.Duration
	RequestDelayJitter time.Duration
	RequestTimeout     time.Duration

	Client *http.ComicClient
	Logger *logger.Logger
}

// TODO: create function to handle creating options with default values and factories for client and logger
// want to avoid having to avoid malformed options
