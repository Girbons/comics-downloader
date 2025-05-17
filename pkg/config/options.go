package config

import (
	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/http"
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

	Client *http.ComicClient
	Logger *logger.Logger
}
