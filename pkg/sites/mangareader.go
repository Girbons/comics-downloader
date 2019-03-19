package sites

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

func retrieveMangareaderImageLinks(c *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(c.URLSource)

	if err != nil {
		log.WithFields(log.Fields{
			"url": c.URLSource,
		}).Error("[MANGAREADER] Something went wrong", err)
	}

	doc := soup.HTMLParse(response)
	// retrieve the <option> tag
	options := doc.FindAll("option")

	for i := 1; i <= len(options); i++ {
		pageLink := fmt.Sprintf("https://%s/%s/%s/%d", c.Source, c.Name, c.IssueNumber, i)
		rsp, soupErr := soup.Get(pageLink)
		if soupErr != nil {
			log.WithFields(log.Fields{
				"url": pageLink,
			}).Error("[MANGAREADER] Oops something went wrong while parsing the current url", soupErr)
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
func SetupMangaReader(c *core.Comic) error {
	name := c.SplitURL()[3]
	IssueNumber := c.SplitURL()[4]

	c.SetName(name)
	c.SetIssueNumber(IssueNumber)

	links, err := retrieveMangareaderImageLinks(c)
	c.SetImageLinks(links)

	return err
}
