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

func isSingleIssue(url string) bool {
	return util.TrimAndSplitURL(url)[3] != "comic"
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func RetrieveIssueLinks(url string, all bool) ([]string, error) {
	if all && isSingleIssue(url) {
		url = "https://www.comicextra.com/comic/" + util.TrimAndSplitURL(url)[3]
	} else if isSingleIssue(url) {
		return []string{url}, nil
	}

	name := util.TrimAndSplitURL(url)[4]
	var links []string
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

// Initialize will initialize the comic based
// on comicextra.com
func Initialize(comic *core.Comic) error {
	parts := util.TrimAndSplitURL(comic.URLSource)
	comic.Name = parts[3]
	comic.IssueNumber = parts[4]

	links, err := retrieveImageLinks(comic)
	comic.Links = links

	return err
}
