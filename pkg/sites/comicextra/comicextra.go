package comicextra

import (
	"regexp"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

type Comicextra struct{}

func (c *Comicextra) retrieveImageLinks(comic *core.Comic) ([]string, error) {
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

func (c *Comicextra) isSingleIssue(url string) bool {
	return util.TrimAndSplitURL(url)[3] != "comic"
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (c *Comicextra) RetrieveIssueLinks(url string, all bool) ([]string, error) {
	if all && c.isSingleIssue(url) {
		url = "https://www.comicextra.com/comic/" + util.TrimAndSplitURL(url)[3]
	} else if c.isSingleIssue(url) {
		return []string{url}, nil
	}

	name := util.TrimAndSplitURL(url)[4]
	var links []string

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("<a[^>]+href=\"([^\">]+" + "/" + name + "/.+)\"")
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1] + "/full"
		if util.IsURLValid(url) && !util.IsValueInSlice(url, links) {
			links = append(links, url)
		}
	}

	return links, err
}

func (c *Comicextra) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	name := parts[3]
	issueNumber := parts[4]

	return name, issueNumber
}

// Initialize will initialize the comic based
// on comicextra.com
func (c *Comicextra) Initialize(comic *core.Comic) error {
	links, err := c.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
