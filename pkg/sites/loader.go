package sites

import (
	"errors"
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/sites/comicextra"
	"github.com/Girbons/comics-downloader/pkg/sites/mangareader"
	"github.com/Girbons/comics-downloader/pkg/sites/mangarock"
	"github.com/Girbons/comics-downloader/pkg/sites/mangatown"
)

func initializeCollection(initializer func(*core.Comic) error, url string, issues []string, conf *config.ComicConfig, format string, source string, options map[string]string) ([]*core.Comic, error) {
	var collection []*core.Comic
	var err error

	if len(issues) == 0 {
		return collection, errors.New("No issues found.")
	}

	for _, url := range issues {
		comic := &core.Comic{
			URLSource: url,
			Config:    conf,
			Source:    source,
			Format:    format,
			Options:   options,
		}
		if err = initializer(comic); err != nil {
			return collection, err
		}
		collection = append(collection, comic)
	}

	return collection, nil
}

// LoadComicFromSource will return a `comic` instance initialized based on the source
func LoadComicFromSource(conf *config.ComicConfig, source, url, country, format string) ([]*core.Comic, error) {
	var err error
	var collection []*core.Comic
	var initializer func(*core.Comic) error
	var issues []string

	options := map[string]string{"country": country}

	switch source {
	case "www.comicextra.com":
		issues, err = comicextra.RetrieveIssueLinks(url)
		initializer = comicextra.Initialize
	case "mangarock.com":
		issues, err = mangarock.RetrieveIssueLinks(url, options)
		initializer = mangarock.Initialize
	case "www.mangareader.net":
		issues, err = mangareader.RetrieveIssueLinks(url)
		initializer = mangareader.Initialize
	case "www.mangatown.com":
		issues, err = mangatown.RetrieveIssueLinks(url)
		initializer = mangatown.Initialize
	default:
		err = fmt.Errorf("It was not possible to determine the source")
	}

	if err != nil {
		return collection, err
	}

	collection, err = initializeCollection(initializer, url, issues, conf, format, source, options)

	return collection, err
}
