package sites

import (
	"context"
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

// Mangareader represents a Mangareader instance
type Mangareader struct {
	options *config.Options
}

// NewMangareader returns a mangareader instance.
func NewMangareader(options *config.Options) *Mangareader {
	return &Mangareader{
		options: options,
	}
}

func init() {
	SupportedSites["mangareader"] = SupportedSite{
		IsEnabled: true,
		Loader:    func(opts *config.Options) BaseSite { return NewMangareader(opts) },
	}
}

func (m *Mangareader) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), m.options.RequestTimeout)
}

func (m *Mangareader) retrieveImageLinks(comic *core.ComicIssue) ([]string, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	response, err := m.options.Client.FetchHTML(ctx, comic.Source.URL)
	if err != nil {
		return nil, err
	}

	var links []string
	doc := soup.HTMLParse(response)
	for _, t := range doc.FindAll("img") {
		imageURL := t.Attrs()["data-src"]
		if util.IsURLValid(imageURL) {
			links = append(links, imageURL)
		}
	}

	if m.options.Debug {
		m.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

func (m *Mangareader) isSingleIssue(url string) bool {
	return len(util.TrimAndSplitURL(url)) >= 5
}

func (m *Mangareader) retrieveLastIssue(url string) (string, error) {
	url = strings.Join(util.TrimAndSplitURL(url)[:4], "/")

	ctx, cancel := m.requestContext()
	defer cancel()

	response, err := m.options.Client.FetchHTML(ctx, url)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(response)
	lastIssue := doc.Find("ul", "class", "d44").FindAll("li")[0].Find("a").Attrs()["href"]
	lastIssueURL := "https://mangareader.tv" + lastIssue

	return lastIssueURL, nil
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (m *Mangareader) RetrieveIssueLinks() ([]string, error) {
	url := m.options.URL

	if m.options.Last {
		lastIssue, err := m.retrieveLastIssue(url)
		return []string{lastIssue}, err
	}

	if m.options.All && m.isSingleIssue(url) {
		url = strings.Join(util.TrimAndSplitURL(url)[:4], "/")
	} else if m.isSingleIssue(url) {
		return []string{url}, nil
	}

	ctx, cancel := m.requestContext()
	defer cancel()

	response, err := m.options.Client.FetchHTML(ctx, url)
	if err != nil {
		return nil, err
	}

	var links []string
	doc := soup.HTMLParse(response)
	nodes := doc.Find("table", "class", "d48").FindAll("tr")
	for _, node := range nodes {
		element := node.Find("a")
		if !strings.Contains(element.NodeValue, "not found") && element.Pointer != nil {
			url := "https://mangareader.tv" + element.Attrs()["href"]
			if util.IsURLValid(url) {
				links = append(links, url)
			}
		}
	}

	if m.options.Debug {
		m.options.Logger.Debug(fmt.Sprintf("Issues Links found: %s", strings.Join(links, " ")))
	}

	return links, err
}

// GetInfo extracts the basic info from the given URL.
func (m *Mangareader) GetInfo(url string) (string, string, error) {
	parts := util.TrimAndSplitURL(url)
	name := parts[3]
	issueNumber := parts[4]

	return name, issueNumber, nil
}

// Initialize loads links and metadata from mangareader
func (m *Mangareader) Initialize(comic *core.ComicIssue) error {
	name, issueNumber, err := m.GetInfo(comic.Source.URL)
	if err != nil {
		return err
	}

	comic.Name = name
	comic.IssueNumber = issueNumber

	links, err := m.retrieveImageLinks(comic)
	comic.ImageLinks = links

	return err
}
