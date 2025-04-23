package sites

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

var baseUrl = "https://readcomiconline.li"

// ReadComicOnline represents a readcomiconline instance.
type ReadComicOnline struct {
	options *config.Options
}

// NewReadComiconline returns a readcomiconline instance.
func NewReadComiconline(options *config.Options) *ReadComicOnline {
	return &ReadComicOnline{
		options: options,
	}
}

func deobfuscateUrl(imageLink string) (string, error) {
	imageLink = strings.ReplaceAll(imageLink, "_x236", "d")
	imageLink = strings.ReplaceAll(imageLink, "_x945", "g")

	if strings.HasPrefix(imageLink, "https://2.bp.blogspot.com") {
		return imageLink, nil
	}

	var quality string

	if strings.Contains(imageLink, "=s0?") {
		imageLink = imageLink[:strings.Index(imageLink, "=s0?")]
		quality = "=s0"
	} else {
		imageLink = imageLink[:strings.Index(imageLink, "=s1600?")]
		quality = "=s1600"
	}

	imageLink = imageLink[4:22] + imageLink[25:]
	imageLink = imageLink[0:len(imageLink)-6] + imageLink[len(imageLink)-2:]

	sd, err := base64.StdEncoding.DecodeString(imageLink)
	if err != nil {
		return "", err
	}

	imageLink = string(sd)
	imageLink = imageLink[0:13] + imageLink[17:]
	imageLink = imageLink[0 : len(imageLink)-2]
	imageLink = imageLink + quality

	link := "https://2.bp.blogspot.com/" + imageLink
	return link, nil
}

func (c *ReadComicOnline) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string

	comic.URLSource = strings.Split(comic.URLSource, "?")[0]

	response, err := soup.Get(comic.URLSource + "?quality=hd&readType=1")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`push\(\'(.*?)\'\)`)
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1]

		clearUrl, err := deobfuscateUrl(url)
		if err != nil {
			return links, err
		}

		if util.IsURLValid(clearUrl) {
			links = append(links, clearUrl)
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
	lastIssue = baseUrl + strings.Split(match[0][1], "?")[0]

	return lastIssue, nil
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (c *ReadComicOnline) RetrieveIssueLinks() ([]string, error) {
	url := c.options.URL

	if c.options.Last {
		issue, err := c.retrieveLastIssue(url)
		return []string{issue}, err
	}

	if c.options.All && c.isSingleIssue(url) {
		url = baseUrl + "/Comic/" + util.TrimAndSplitURL(url)[3]
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
			url = baseUrl + strings.Split(url, "?")[0]
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

// GetInfo extracts the basic info from the given url.
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
