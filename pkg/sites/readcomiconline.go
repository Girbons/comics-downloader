package sites

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

type ReadComicOnline struct {
	options *config.Options
}

func NewReadComiconline(options *config.Options) *ReadComicOnline {
	return &ReadComicOnline{
		options: options,
	}
}

func (c *ReadComicOnline) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string

	comic.URLSource = strings.Split(comic.URLSource, "?")[0]

	response, err := soup.Get(comic.URLSource + "?quality=hd&readType=1")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`push\(\"(.*?)\"\)`)
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

func (c *ReadComicOnline) isSingleIssue(url string) bool {
	parts := util.TrimAndSplitURL(url)
	return len(parts) > 5 && strings.Contains(parts[5], "Issue-")
}

func (c *ReadComicOnline) retrieveLastIssue(url string) (string, error) {
	var lastIssue string

	response, err := soup.Get(url)
	if err != nil {
		return "", err
	}

	name := util.TrimAndSplitURL(url)[4]
	re := regexp.MustCompile("<a[^>]+href=\"([^\">]+" + "/" + name + "/.+)\"")
	match := re.FindAllStringSubmatch(response, -1)
	lastIssue = "https://readcomiconline.to" + strings.Split(match[0][1], "?")[0]

	return lastIssue, nil
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (c *ReadComicOnline) RetrieveIssueLinks() ([]string, error) {
	url := c.options.Url

	if c.options.Last {
		issue, err := c.retrieveLastIssue(url)
		return []string{issue}, err
	}

	if c.options.All && c.isSingleIssue(url) {
		url = "https://readcomiconline.to/Comic/" + util.TrimAndSplitURL(url)[3]
	} else if c.isSingleIssue(url) {
		return []string{url}, nil
	}

	name := util.TrimAndSplitURL(url)[4]
	var (
		pages []string
		links []string
	)

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	pages = append(pages, url)
	re := regexp.MustCompile("<a[^>]+href=\"([^\">]+" + "/" + name + "/.+)\"")
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1]
		if !util.IsValueInSlice(url, pages) {
			url = "https://readcomiconline.to" + strings.Split(url, "?")[0]
			if util.IsURLValid(url) && !util.IsValueInSlice(url, links) {
				links = append(links, url)
			}
		}
	}

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Issues Links retrieved: %s", strings.Join(links, " ")))
	}

	return links, err
}

func (c *ReadComicOnline) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	name := parts[4]
	issueNumber := strings.Split(strings.Replace(parts[5], "Issue-", "", -1), "?")[0]

	return name, issueNumber
}

// Initialize will initialize the comic based
// on ReadComicOnline.to
func (c *ReadComicOnline) Initialize(comic *core.Comic) error {
	links, err := c.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
