package sites

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

// Comicextra represents comicextra instance.
type Comicextra struct {
	options *config.Options
}

// NewComicextra returs a comicextra instance.
func NewComicextra(options *config.Options) *Comicextra {
	return &Comicextra{
		options: options,
	}
}

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

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

func (c *Comicextra) isSingleIssue(url string) bool {
	return util.TrimAndSplitURL(url)[3] != "comic"
}

func (c *Comicextra) retrieveLastIssue(url string) (string, error) {
	var lastIssue string

	response, err := soup.Get(url)

	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(response)

	issues := doc.FindAll("option")

	var validLinks []string

	for _, v := range issues {
		tmpUrl := v.Attrs()["value"]
		if util.IsURLValid(tmpUrl) {
			validLinks = append(validLinks, tmpUrl)
		}
	}

	sort.Strings(validLinks)

	lastIssue = validLinks[len(validLinks)-1]

	return lastIssue, nil
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (c *Comicextra) RetrieveIssueLinks() ([]string, error) {
	url := c.options.URL

	splittedUrl := util.TrimAndSplitURL(url)

	comicName := splittedUrl[3]

	if c.options.Last {
		issue, err := c.retrieveLastIssue(url)
		return []string{issue}, err
	}

	// check if the user has submitted a generic url
	if !strings.HasSuffix(url, "/issue") {
		// enable `-all` flag in this case
		c.options.Logger.Warning("URL does not contain a specific issue, `-all` flag will be automaticcaly enabled")
		c.options.All = true
	}

	if c.options.All && c.isSingleIssue(url) {
		url = "https://" + c.options.Source + "/comic/" + comicName
	} else if c.isSingleIssue(url) {

		if !strings.HasSuffix(url, "/full") {
			url = url + "/full"
		}

		return []string{url}, nil
	}

	var (
		links    []string
		elements []soup.Root
	)

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)

	episodeList := doc.Find("div", "class", "episode-list")
	if episodeList.Pointer != nil {
		elements = episodeList.FindAll("a")
	}

	for _, element := range elements {
		pageURL := element.Attrs()["href"]
		if !strings.Contains(pageURL, "onclick") && !util.IsValueInSlice(pageURL, links) {
			if !strings.HasSuffix(pageURL, "/full") {
				pageURL = pageURL + "/full"
			}

			links = append(links, pageURL)
		}
	}

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Issues Links retrieved: %s", strings.Join(links, " ")))
	}

	return links, err
}

// GetInfo extracts the basic info from the given url.
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
