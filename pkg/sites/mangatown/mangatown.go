package mangatown

import (
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/core"
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

// SetupMangaTown will initialize the comic based
// on mangatown.com
func Initialize(comic *core.Comic) error {
	comic.Name = comic.SplitURL()[4]
	comic.IssueNumber = comic.SplitURL()[6]

	links, err := retrieveImageLinks(comic)
	comic.Links = links

	return err
}
