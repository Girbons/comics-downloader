package sites

import (
	"regexp"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

func retrieveComicExtraImageLinks(c *core.Comic) ([]string, error) {
	response, err := soup.Get(c.URLSource)

	if err != nil {
		log.Error("[COMICEXTRA] Something went wrong with: ", c.URLSource, err)
	}

	re := regexp.MustCompile(c.ImageRegex)
	match := re.FindAllStringSubmatch(response, -1)

	links := make([]string, len(match))

	for i := range links {
		url := match[i][1]
		if util.IsUrlValid(url) {
			links[i] = url
		}
	}

	return links, err

}

// SetupComicExtra will initialize the comic based
// on comicextra.com
func SetupComicExtra(c *core.Comic) error {
	name := c.SplitURL()[3]
	issueNumber := c.SplitURL()[4]
	imageRegex := `<img[^>]+src="([^">]+)"`
	c.SetInfo(name, issueNumber, imageRegex)

	links, err := retrieveComicExtraImageLinks(c)
	c.SetImageLinks(links)

	return err
}
