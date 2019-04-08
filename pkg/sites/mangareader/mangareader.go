package mangareader

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/anaskhan96/soup"
)

func retrieveImageLinks(c *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(c.URLSource)

	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(response)
	// retrieve the <option> tag
	options := doc.FindAll("option")

	for i := 1; i <= len(options); i++ {
		pageLink := fmt.Sprintf("https://%s/%s/%s/%d", c.Source, c.Name, c.IssueNumber, i)
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

// SetupMangaReader will initialize the comic based
// www.mangareader.net
func Initialize(c *core.Comic) error {
	name := c.SplitURL()[3]
	IssueNumber := c.SplitURL()[4]

	c.SetName(name)
	c.SetIssueNumber(IssueNumber)

	links, err := retrieveImageLinks(c)
	c.SetImageLinks(links)

	return err
}
