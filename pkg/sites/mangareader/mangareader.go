package mangareader

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/anaskhan96/soup"
)

func retrieveImageLinks(comic *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(comic.URLSource)

	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)
	// retrieve the <option> tag
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

// Initialize loads links and metadata from mangareader
func Initialize(comic *core.Comic) error {
	comic.Name = comic.SplitURL()[3]
	comic.IssueNumber = comic.SplitURL()[4]

	links, err := retrieveImageLinks(comic)
	comic.Links = links

	return err
}
