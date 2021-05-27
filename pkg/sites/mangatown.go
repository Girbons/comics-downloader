package sites

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

// Mangatown represents a Mangatown instance.
type Mangatown struct {
	options *config.Options
}

// NewMangatown returns a new mangatown instance.
func NewMangatown(options *config.Options) *Mangatown {
	return &Mangatown{
		options: options,
	}
}

func (m *Mangatown) findPages(document *soup.Root) []string {
	var pages []string

	options := document.Find("div", "class", "page_select").Find("select").FindAll("option")

	for _, option := range options {
		if !strings.Contains(option.Text(), "Featured") {
			pages = append(pages, option.Text())
		}
	}

	return pages
}

func (m *Mangatown) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string
	var link string

	response, err := soup.Get(comic.URLSource)

	if err != nil {
		return nil, err
	}

	document := soup.HTMLParse(response)
	pages := m.findPages(&document)

	for _, page := range pages {
		link = fmt.Sprintf("%s%s.html", comic.URLSource, page)
		response, err := soup.Get(link)

		if err != nil {
			return nil, err
		}

		document = soup.HTMLParse(response)
		img := document.Find("div", "id", "viewer").Find("a").Find("img")
		links = append(links, fmt.Sprintf("https:%s", img.Attrs()["src"]))
	}

	if m.options.Debug {
		m.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

func (m *Mangatown) isSingleIssue(url string) bool {
	parts := util.TrimAndSplitURL(url)
	return len(parts) >= 6 && parts[5] != ""
}

func (m *Mangatown) retrieveLastIssue(url string) (string, error) {
	url = strings.Join(util.TrimAndSplitURL(url)[:5], "/")
	response, err := soup.Get(url)

	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(response)
	chapters := doc.Find("ul", "class", "chapter_list").FindAll("a")

	// the first element is the last chapter
	lastIssue := "https://www.mangatown.com" + chapters[0].Attrs()["href"]
	return lastIssue, nil
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (m *Mangatown) RetrieveIssueLinks() ([]string, error) {
	url := m.options.URL

	if m.options.Last {
		lastIssue, err := m.retrieveLastIssue(url)
		return []string{lastIssue}, err
	}

	if m.options.All && m.isSingleIssue(url) {
		url = strings.Join(util.TrimAndSplitURL(url)[:5], "/")
	} else if m.isSingleIssue(url) {
		return []string{url}, nil
	}

	var links []string

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)
	chapters := doc.Find("ul", "class", "chapter_list").FindAll("a")

	for _, chapter := range chapters {
		url := "https://mangatown.com" + chapter.Attrs()["href"]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	if m.options.Debug {
		m.options.Logger.Debug(fmt.Sprintf("Issue Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

// GetInfo extracts the basic info from the given URL.
func (m *Mangatown) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	name := parts[4]
	issueNumber := parts[len(parts)-1]

	return name, issueNumber
}

// Initialize loads links and metadata from mangatown
func (m *Mangatown) Initialize(comic *core.Comic) error {
	links, err := m.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
