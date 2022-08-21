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

var baseUrl = "https://readcomiconline.li"
var highQuality = "?quality=hd&readType=1"

// ReadComicOnline represents a readcomiconline instance.
type ReadComicOnline struct {
	options *config.Options
}

// NewReadComiconline returns a readcomiconline instance.
func NewReadComiconline(options *config.Options) *ReadComicOnline {
	return &ReadComicOnline{options: options}
}

func (c *ReadComicOnline) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var (
		links []string
		pages []string
	)

	comic.URLSource = strings.Split(comic.URLSource, "?")[0]

	response, err := soup.Get(comic.URLSource + highQuality)
	if err != nil {
		return nil, err
	}

	document := soup.HTMLParse(response)

	if strings.HasSuffix(c.options.Source, ".ru") {
		results := document.Find("select", "id", "page-list")
		// extract pages
		for _, el := range results.FindAll("option") {
			pages = append(pages, el.Attrs()["value"])
		}

		for _, page := range pages {
			url := comic.URLSource

			if !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += page + highQuality

			resp, _ := soup.Get(url)
			inner_document := soup.HTMLParse(resp)

			for _, l := range inner_document.FindAll("img") {
				image_link := strings.Replace(l.Attrs()["src"], " ", "", -1)

				if image_link != "" && strings.Contains(image_link, "chapters") {
					links = append(links, image_link)
				}
			}

		}
	} else {
		fmt.Println(response)
		re := regexp.MustCompile(`push\(\'(.*?)\'\)`)
		match := re.FindAllStringSubmatch(response, -1)

		baseImageUrl := "https://2.bp.blogspot.com/"

		for i := range match {
			url := baseImageUrl + match[i][1]
			if util.IsURLValid(url) {
				links = append(links, url)
			}
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
	var links []string

	url := c.options.URL

	if c.options.Last {
		issue, err := c.retrieveLastIssue(url)
		return []string{issue}, err
	}

	if c.options.All {
		url = "https://" + c.options.Source + "/Comic/" + util.TrimAndSplitURL(url)[4]
	}

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	name := util.TrimAndSplitURL(url)[4]

	re := regexp.MustCompile("<a[^>]+href=\"([^\">]+" + "/" + name + "/.+)\"")
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1]

		if !strings.HasPrefix(url, ".ru") {
			url = baseUrl + strings.Split(url, "?")[0]
		}

		if util.IsURLValid(url) && !util.IsValueInSlice(url, links) {
			links = append(links, url)
		}
	}

	if c.options.Debug {
		c.options.Logger.Debug(fmt.Sprintf("Issues Links retrieved: %s", strings.Join(links, " ")))
	}

	return links, nil
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
