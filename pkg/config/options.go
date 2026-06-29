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
	Format              string
	CustomComicName     string
	ForceAspect         bool
	OutputFolder        string
	CreateDefaultPath   bool
	IssueNumberNameOnly bool
	URL                 string
	Source              string
	IssuesRange         string
	IssueFolderName     string
	UserAgents          []string
	SessionCookie       string
	RequestDelay        time.Duration
	RequestDelayJitter  time.Duration

	Client *http.ComicClient
	Logger *logger.Logger
}
