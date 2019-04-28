package comicextra

import (
	"regexp"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

func retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(comic.URLSource)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(util.IMAGEREGEX)
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	return links, err
}

func RetrieveIssueLinks(url string) ([]string, error) {
	if isSingleIssue(url) {
		return []string{url}, nil
	}

	var links []string
	name := util.SplitURL(url)[4]
	set := make(map[string]struct{})

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("<a[^>]+href=\"([^\">]+" + "/" + name + "/.+)\"")
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1] + "/full"
		if util.IsURLValid(url) {
			set[url] = struct{}{}
		}
	}

	for url := range set {
		links = append(links, url)
	}

	return links, err
}

func isSingleIssue(url string) bool {
	return util.SplitURL(url)[3] != "comic"
}

// Initialize will initialize the comic based
// on comicextra.com
func Initialize(comic *core.Comic) error {
	comic.Name = comic.SplitURL()[3]
	comic.IssueNumber = comic.SplitURL()[4]

	links, err := retrieveImageLinks(comic)
	comic.Links = links

	return err
}
