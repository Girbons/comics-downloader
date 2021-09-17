package config

import "github.com/Girbons/comics-downloader/internal/logger"

// Options represents the comics downloader options.
type Options struct {
	Debug             bool
	All               bool
	Last              bool
	ImagesOnly        bool
	Daemon            bool
	ImagesFormat      string
	Country           string
	Format            string
	ForceAspect       bool
	OutputFolder      string
	CreateDefaultPath bool
	URL               string
	Source            string
	IssuesRange       string
	Timeout           int

	Logger *logger.Logger
}
