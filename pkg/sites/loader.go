package sites

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Girbons/comics-downloader/internal/flag/parser"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
)

func initializeCollection(issues []string, options *config.Options, base BaseSite) ([]*core.Comic, error) {
	var collection []*core.Comic
	var err error

	if len(issues) == 0 {
		return collection, errors.New("No issues found")
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

		dir, _ := util.PathSetup(options.CreateDefaultPath, options.OutputFolder, options.Source, name)
		fileName := util.GetPathToFile(dir, name, issueNumber, options.Format, options.IssueNumberNameOnly)

		if util.DirectoryOrFileDoesNotExist(fileName) || options.ImagesOnly {
			comic := &core.Comic{
				Name:         name,
				IssueNumber:  issueNumber,
				URLSource:    url,
				Source:       options.Source,
				Format:       options.Format,
				ImagesFormat: options.ImagesFormat,
			}
			if err = base.Initialize(comic); err != nil {
				return collection, err
			}
			collection = append(collection, comic)
		}
	}

	return collection, nil
}

var onlyNumbers = regexp.MustCompile("[^0-9]+[^.][^0-9]+")

func notInIssuesRange(issueNumber string, start, end float64) bool {
	if start == 0 || end == 0 {
		return false
	}

	normalizedNumber := onlyNumbers.ReplaceAllString(issueNumber, "")
	if normalizedNumber == "" {
		return true
	}

	number, err := strconv.ParseFloat(normalizedNumber, 64)
	if err != nil {
		return true
	}

	return number < start || number > end
}

// LoadComicFromSource will return an `comic` instance initialized based on the source
func LoadComicFromSource(options *config.Options) ([]*core.Comic, error) {
	var (
		base       BaseSite
		issues     []string
		collection []*core.Comic
		err        error
	)

	switch {
	case strings.Contains(options.Source, "readcomiconline"):
		base = NewReadComiconline(options)
	case strings.Contains(options.Source, "comicextra"):
		base = NewComicextra(options)
	case strings.Contains(options.Source, "mangareader"):
		base = NewMangareader(options)
	case strings.Contains(options.Source, "mangatown"):
		base = NewMangatown(options)
	case strings.Contains(options.Source, "mangadex"):
		base = NewMangadex(options)
	case strings.Contains(options.Source, "readallcomics"):
		base = NewReadallcomics(options)
	case strings.Contains(options.Source, "mangakakalot"):
		base = NewMangaKakalot(options)
	case strings.Contains(options.Source, "manganato"):
		base = NewManganato(options)
	default:
		err = fmt.Errorf("source unknown")
		return collection, err
	}

	issues, err = base.RetrieveIssueLinks()
	if err != nil {
		return collection, err
	}

	return initializeCollection(issues, options, base)
}
