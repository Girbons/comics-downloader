package sites

import (
	"errors"
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
)

func initializeCollection(issues []string, url, format, source, imagesFormat string, siteSource BaseSite, imagesOnly bool, outputFolder string) ([]*core.Comic, error) {
	var collection []*core.Comic
	var err error

	if len(issues) == 0 {
		return collection, errors.New("No issues found.")
	}

	for _, url := range issues {
		name, issueNumber := GetInfo(siteSource, url)
		name = util.Parse(name)
		issueNumber = util.Parse(issueNumber)

		dir, _ := util.PathSetup(outputFolder, source, name)
		fileName := util.GenerateFileName(dir, name, issueNumber, format)

		if util.DirectoryOrFileDoesNotExist(fileName) || imagesOnly {
			comic := &core.Comic{
				Name:         name,
				IssueNumber:  issueNumber,
				URLSource:    url,
				Source:       source,
				Format:       format,
				ImagesFormat: imagesFormat,
			}
			if err = Initialize(siteSource, comic); err != nil {
				return collection, err
			}
			collection = append(collection, comic)
		}
	}

	return collection, nil
}

// LoadComicFromSource will return a `comic` instance initialized based on the source
func LoadComicFromSource(source, url, country, format, imagesFormat string, all, last, imagesOnly bool, outputFolder string) ([]*core.Comic, error) {
	var siteSource BaseSite
	var collection []*core.Comic
	var issues []string
	var err error

	options := map[string]string{"country": country}

	switch source {
	case "www.comicextra.com":
		siteSource = &Comicextra{}
	case "mangarock.com":
		siteSource = NewMangarock(options)
	case "www.mangareader.net":
		siteSource = &Mangareader{}
	case "www.mangatown.com":
		siteSource = &Mangatown{}
	case "mangadex.cc":
		siteSource = NewMangadex(country)
	default:
		err = fmt.Errorf("It was not possible to determine the source")
		return collection, err
	}

	issues, err = RetrieveIssueLinks(siteSource, url, all, last)
	if err != nil {
		return collection, err
	}

	return initializeCollection(issues, url, format, source, imagesFormat, siteSource, imagesOnly, outputFolder)
}
