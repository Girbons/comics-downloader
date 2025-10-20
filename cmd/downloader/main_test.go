package main

import (
	"testing"
)

func TestBuildOptionsCopiesGlobals(t *testing.T) {
	prev := struct {
		debug               bool
		all                 bool
		last                bool
		imagesOnly          bool
		imagesFormat        string
		country             string
		forceAspect         bool
		format              string
		customComicName     string
		issueNumberNameOnly bool
		url                 string
		outputFolder        string
		createDefaultPath   bool
		daemon              bool
		daemonTimeout       int
		issuesRange         string
		issueFolderName     string
	}{
		debug:               debug,
		all:                 all,
		last:                last,
		imagesOnly:          imagesOnly,
		imagesFormat:        imagesFormat,
		country:             country,
		forceAspect:         forceAspect,
		format:              format,
		customComicName:     customComicName,
		issueNumberNameOnly: issueNumberNameOnly,
		url:                 url,
		outputFolder:        outputFolder,
		createDefaultPath:   createDefaultPath,
		daemon:              daemon,
		daemonTimeout:       daemonTimeout,
		issuesRange:         issuesRange,
		issueFolderName:     issueFolderName,
	}
	defer func() {
		debug = prev.debug
		all = prev.all
		last = prev.last
		imagesOnly = prev.imagesOnly
		imagesFormat = prev.imagesFormat
		country = prev.country
		forceAspect = prev.forceAspect
		format = prev.format
		customComicName = prev.customComicName
		issueNumberNameOnly = prev.issueNumberNameOnly
		url = prev.url
		outputFolder = prev.outputFolder
		createDefaultPath = prev.createDefaultPath
		daemon = prev.daemon
		daemonTimeout = prev.daemonTimeout
		issuesRange = prev.issuesRange
		issueFolderName = prev.issueFolderName
	}()

	debug = true
	all = true
	last = true
	imagesOnly = true
	imagesFormat = "png"
	country = "jp"
	forceAspect = true
	format = "epub"
	customComicName = "custom"
	issueNumberNameOnly = true
	url = "http://example.com/comic"
	outputFolder = "/tmp/output"
	createDefaultPath = false
	daemon = true
	daemonTimeout = 42
	issuesRange = "1-5"
	issueFolderName = "chapter-"

	opts := buildOptions()

	if !opts.Debug || !opts.All || !opts.Last || !opts.ImagesOnly {
		t.Fatalf("expected boolean flags to be copied into options: %+v", opts)
	}

	if opts.ImagesFormat != "png" || opts.Country != "jp" || opts.Format != "epub" {
		t.Fatalf("expected string values to be copied, got %+v", opts)
	}

	if opts.CustomComicName != "custom" || opts.URL != "http://example.com/comic" {
		t.Fatalf("expected URL and custom name to be copied, got %+v", opts)
	}

	if opts.OutputFolder != "/tmp/output" || opts.IssueFolderName != "chapter-" {
		t.Fatalf("expected folder values to be copied, got %+v", opts)
	}

	if opts.DaemonTimeout != 42 || !opts.Daemon || opts.CreateDefaultPath {
		t.Fatalf("expected daemon configuration to be copied, got %+v", opts)
	}

	if opts.IssuesRange != "1-5" {
		t.Fatalf("expected issues range to be copied, got %q", opts.IssuesRange)
	}
}
