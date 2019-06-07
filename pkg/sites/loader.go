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
	"github.com/Girbons/comics-downloader/pkg/util"
	log "github.com/sirupsen/logrus"
)

func initializeCollection(issues []string, conf *config.ComicConfig, url, format, source string, options map[string]string, siteLoader *SiteLoader) ([]*core.Comic, error) {
	var collection []*core.Comic
	var err error

	if len(issues) == 0 {
		return collection, errors.New("No issues found.")
	}

	for _, url := range issues {
		name, issueNumber := GetInfo(siteLoader.Source, url, options)
		name = util.Parse(name)
		issueNumber = util.Parse(issueNumber)

		dir, _ := util.PathSetup(source, name)
		fileName := util.GenerateFileName(dir, issueNumber, format)

		if util.FileDoesNotExist(fileName) {
			comic := &core.Comic{
				Name:        name,
				IssueNumber: issueNumber,
				URLSource:   url,
				Config:      conf,
				Source:      source,
				Format:      format,
				Options:     options,
			}
			if err = Initialize(siteLoader.Source, comic); err != nil {
				return collection, err
			}
			collection = append(collection, comic)
		} else {
			log.Info(fmt.Sprintf("%s/%s.%s Already exist", name, issueNumber, format))
		}
	}

	return collection, nil
}

// LoadComicFromSource will return a `comic` instance initialized based on the source
func LoadComicFromSource(conf *config.ComicConfig, source, url, country, format string, all bool) ([]*core.Comic, error) {
	var collection []*core.Comic
	var issues []string
	var err error
	options := map[string]string{"country": country}

	siteLoader := &SiteLoader{}

	switch source {
	case "www.comicextra.com":
		siteLoader.Source = &comicextra.Comicextra{}
	case "mangarock.com":
		siteLoader.Source = &mangarock.Mangarock{}
	case "www.mangareader.net":
		siteLoader.Source = &mangareader.Mangareader{}
	case "www.mangatown.com":
		siteLoader.Source = &mangatown.Mangatown{}
	default:
		err = fmt.Errorf("It was not possible to determine the source")
	}

	issues, err = RetrieveIssueLinks(siteLoader.Source, url, all, options)
	if err != nil {
		return collection, err
	}

	return initializeCollection(issues, conf, url, format, source, options, siteLoader)
}
