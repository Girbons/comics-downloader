package sites

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Girbons/comics-downloader/internal/flag/parser"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
)

func initializeCollection(issues []string, options *config.Options, base BaseSite) ([]*core.Comic, error) {
	var collection []*core.Comic
	var err error

	if len(issues) == 0 {
		return collection, errors.New("No issues found.")
	}

	var startRange, endRange int
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
		issueNumber = util.Parse(issueNumber)

		if notInIssuesRange(issueNumber, startRange, endRange) {
			continue
		}

		dir, _ := util.PathSetup(options.OutputFolder, options.Source, name)
		fileName := util.GenerateFileName(dir, name, issueNumber, options.Format)

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

var onlyNumbers = regexp.MustCompile("[^0-9]+")

func notInIssuesRange(issueNumber string, start, end int) bool {
	if start == 0 || end == 0 {
		return false
	}

	normalizedNumber := onlyNumbers.ReplaceAllString(issueNumber, "")
	if normalizedNumber == "" {
		return true
	}

	number, err := strconv.Atoi(normalizedNumber)
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

	switch options.Source {
	case "readcomiconline.li":
		base = NewReadComiconline(options)
	case "www.comicextra.com":
		base = NewComicextra(options)
	case "www.mangareader.net":
		base = NewMangareader(options)
	case "www.mangatown.com":
		base = NewMangatown(options)
	case "mangadex.cc", "mangadex.org":
		base = NewMangadex(options)
	case "readallcomics.com":
		base = NewReadallcomics(options)
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
