package config

import "github.com/Girbons/comics-downloader/internal/logger"

type Options struct {
	Debug        bool
	All          bool
	Last         bool
	ImagesOnly   bool
	Daemon       bool
	ImagesFormat string
	Country      string
	Format       string
	OutputFolder string
	Url          string
	Source       string
	IssuesRange  string
	Timeout      int

	Logger *logger.Logger
}
