package sites

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Girbons/comics-downloader/internal/flag/parser"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
	"github.com/Girbons/comics-downloader/pkg/util"
)

func initializeCollection(issues []string, options *config.Options, base BaseSite) ([]*core.ComicIssue, error) {
	var collection []*core.ComicIssue
	// var err error

	if len(issues) == 0 {
		return collection, fmt.Errorf("no issues found for URL %q; ensure it points to a specific comic or chapter page", options.URL)
	}

	var startRange, endRange float64
	if options.All && options.IssuesRange != "" {
		start, end, err := parser.ParseIssuesRange(options.IssuesRange)
		if err != nil {
			return collection, err
		}
		startRange = start
		endRange = end
	}

	for _, url := range issues {
		name, issueNumber := base.GetInfo(url)
		name = util.Parse(name)
		if len(options.CustomComicName) > 0 {
			name = options.CustomComicName
		}
		issueNumber = util.Parse(issueNumber)

		if notInIssuesRange(issueNumber, startRange, endRange) {
			continue
		}

		outputFormat, err := core.ToComicOutputFormat(options.OutputFormat)
		if err != nil {
			return collection, err
		}

		dir, pathErr := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, options.SourceName, name)
		if pathErr != nil {
			return collection, pathErr
		}
		fileName := util.GetPathToFile(dir, name, issueNumber, outputFormat.String(), options.IssueNumberNameOnly)

		if util.DirectoryOrFileDoesNotExist(fileName) || options.ImagesOnly {

			comic := &core.ComicIssue{
				Name:        name,
				IssueNumber: issueNumber,

				OutputFormat: outputFormat,
				ImagesFormat: options.ImagesFormat,

				Source: &core.ComicSource{
					Name: options.SourceName,
					URL:  url,
				},
				SeriesMetadata: &core.SeriesMetadata{
					Title: name,
				},
			}
			if err = base.Initialize(comic); err != nil {
				return collection, err
			}
			collection = append(collection, comic)
		}
	}

	return collection, nil
}

var volumeAndIssuePattern = regexp.MustCompile(`v(\d+)[-_]0*(\d+)`)
var onlyDigits = regexp.MustCompile(`\d+`)

func notInIssuesRange(issueNumber string, start, end float64) bool {
	if start == 0 || end == 0 {
		return false
	}

	number := extractIssueNumberForRange(issueNumber)
	if number == 0 {
		return true
	}

	return number < start || number > end
}

// extractIssueNumberForRange extracts a numeric value from an issue number string
// for range comparison. It supports:
//  1. Volume.Issue format where v<vol>-<issue> becomes <vol>.<issue> as a decimal
//     (e.g., "v4-078" -> 4.78 for Volume 4, Issue 78)
//  2. Simple numeric format (e.g., "078" -> 78, "20.5" -> 20.5)
func extractIssueNumberForRange(issueNumber string) float64 {
	// Try to match volume and issue pattern (e.g., "v4-078-2016")
	if matches := volumeAndIssuePattern.FindStringSubmatch(issueNumber); matches != nil {
		volume, _ := strconv.Atoi(matches[1])
		issue, _ := strconv.Atoi(matches[2])
		// Combine as volume.issue decimal (e.g., volume 4, issue 78 becomes 4.78)
		// This allows users to specify ranges like "4.78-4.99" for V4 issues 78-99
		return float64(volume) + float64(issue)/100.0
	}

	// Try to parse as a simple float (e.g., "20.5")
	if number, err := strconv.ParseFloat(issueNumber, 64); err == nil {
		return number
	}

	// Extract the first sequence of digits (e.g., "078" from "078-something")
	if matches := onlyDigits.FindString(issueNumber); matches != "" {
		if number, err := strconv.ParseFloat(matches, 64); err == nil {
			return number
		}
	}

	return 0
}

// LoadComicFromSource will return an `comic` instance initialized based on the source
func LoadComicFromSource(options *config.Options) ([]*core.ComicIssue, error) {
	var (
		base       BaseSite
		issues     []string
		collection []*core.ComicIssue
		err        error
	)

	// ensure the client is actually set
	if options.Client == nil {
		options.Client = httpclient.NewComicClient()
	}

	switch {
	case strings.Contains(options.SourceName, "readcomiconline"):
		base = NewReadComiconline(options)
	case strings.Contains(options.SourceName, "comicextra"):
		base = NewComicextra(options)
	case strings.Contains(options.SourceName, "mangareader"):
		base = NewMangareader(options)
	case strings.Contains(options.SourceName, "mangatown"):
		base = NewMangatown(options)
	case strings.Contains(options.SourceName, "mangadex"):
		base = NewMangadex(options)
	case strings.Contains(options.SourceName, "readallcomics"):
		base = NewReadallcomics(options)
	case strings.Contains(options.SourceName, "mangakakalot"):
		base = NewMangaKakalot(options)
	case strings.Contains(options.SourceName, "manganato"):
		base = NewManganato(options)
	default:
		err = fmt.Errorf("source unknown")
		return collection, err
	}

	if options.Logger != nil && options.Debug {
		options.Logger.Debugf("sites: retrieving issues for %s", options.URL)
	}

	issues, err = base.RetrieveIssueLinks()
	if err != nil {
		return collection, err
	}

	if options.Logger != nil && options.Debug {
		options.Logger.Debugf("sites: %d issue(s) discovered for %s", len(issues), options.URL)
	}

	return initializeCollection(issues, options, base)
}
