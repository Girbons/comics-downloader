package mangatown

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
)

func findPages(document *soup.Root) []string {
	var pages []string

	options := document.Find("div", "class", "page_select").Find("select").FindAll("option")

	for _, option := range options {
		if !strings.Contains(option.Text(), "Featured") {
			pages = append(pages, option.Text())
		}
	}

	return pages
}

func retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string
	var link string

	response, err := soup.Get(comic.URLSource)

	if err != nil {
		return nil, err
	}

	document := soup.HTMLParse(response)
	pages := findPages(&document)

	for _, page := range pages {
		link = fmt.Sprintf("%s%s.html", comic.URLSource, page)
		response, err := soup.Get(link)

		if err != nil {
			return nil, err
		}

		document = soup.HTMLParse(response)
		img := document.Find("div", "id", "viewer").Find("a").Find("img")
		links = append(links, img.Attrs()["src"])
	}

	return links, err
}

func isSingleIssue(url string) bool {
	parts := util.TrimAndSplitURL(url)
	return len(parts) >= 6 && parts[5] != ""
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func RetrieveIssueLinks(url string) ([]string, error) {
	if isSingleIssue(url) {
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
		url := "https:" + chapter.Attrs()["href"]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	return links, err
}

// Initialize loads links and metadata from mangatown
func Initialize(comic *core.Comic) error {
	parts := util.TrimAndSplitURL(comic.URLSource)
	comic.Name = parts[4]
	comic.IssueNumber = parts[len(parts)-1]

	links, err := retrieveImageLinks(comic)
	comic.Links = links

	return err
}
