package sites

import (
	"errors"
	"fmt"

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

	for _, url := range issues {
		name, issueNumber := base.GetInfo(url)
		name = util.Parse(name)
		issueNumber = util.Parse(issueNumber)

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

// LoadComicFromSource will return an `comic` instance initialized based on the source
func LoadComicFromSource(options *config.Options) ([]*core.Comic, error) {
	var (
		base       BaseSite
		issues     []string
		collection []*core.Comic
		err        error
	)

	switch options.Source {
	case "readcomiconline.to":
		base = &ReadComicOnline{}
	case "www.comicextra.com":
		base = &Comicextra{}
	case "www.mangareader.net":
		base = &Mangareader{}
	case "www.mangatown.com":
		base = &Mangatown{}
	case "mangadex.cc", "mangadex.org":
		base = NewMangadex(options.Country, options.Source)
	default:
		err = fmt.Errorf("source unknown")
		return collection, err
	}

	issues, err = base.RetrieveIssueLinks(options.Url, options.All, options.Last)
	if err != nil {
		return collection, err
	}

	return initializeCollection(issues, options, base)
}
