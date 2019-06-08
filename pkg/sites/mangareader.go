package sites

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

type Mangareader struct{}

func (m *Mangareader) retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(comic.URLSource)

	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)
	// retrieve the <option>
	options := doc.FindAll("option")

	for i := 1; i <= len(options); i++ {
		pageLink := fmt.Sprintf("https://%s/%s/%s/%d", comic.Source, comic.Name, comic.IssueNumber, i)
		rsp, soupErr := soup.Get(pageLink)
		if soupErr != nil {
			return nil, soupErr
		}

		doc = soup.HTMLParse(rsp)
		// return the first `<img>`
		// e.g. <img src="http://example.com">
		imgTag := doc.Find("img")
		// doc.Find returns an html.Node
		// the line below will return the src value
		src := imgTag.Pointer.Attr[3].Val
		links = append(links, src)
	}

	return links, err
}

func (m *Mangareader) isSingleIssue(url string) bool {
	return len(util.TrimAndSplitURL(url)) >= 5
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (m *Mangareader) RetrieveIssueLinks(url string, all bool) ([]string, error) {
	if all && m.isSingleIssue(url) {
		url = strings.Join(util.TrimAndSplitURL(url)[:4], "/")
	} else if m.isSingleIssue(url) {
		return []string{url}, nil
	}

	var links []string

	response, err := soup.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)
	chapters := doc.Find("div", "id", "chapterlist").FindAll("a")

	for _, chapter := range chapters {
		url := "https://www.mangareader.net" + chapter.Attrs()["href"]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	return links, err
}

func (m *Mangareader) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	name := parts[3]
	issueNumber := parts[4]

	return name, issueNumber
}

// Initialize loads links and metadata from mangareader
func (m *Mangareader) Initialize(comic *core.Comic) error {
	name, issueNumber := m.GetInfo(comic.URLSource)
	comic.Name = name
	comic.IssueNumber = issueNumber

	links, err := m.retrieveImageLinks(comic)
	comic.Links = links

	return err
}
