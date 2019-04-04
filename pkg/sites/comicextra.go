package sites

import (
	"regexp"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

func retrieveComicExtraImageLinks(c *core.Comic) ([]string, error) {
	var links []string

	response, err := soup.Get(c.URLSource)
	if err != nil {
		log.WithFields(log.Fields{"source": c.Source}).Error(err)
	}

	re := regexp.MustCompile(util.IMAGEREGEX)
	match := re.FindAllStringSubmatch(response, -1)

	for i := range match {
		url := match[i][1]
		if util.IsURLValid(url) {
			links = append(links, url)
		}
	}

	return links, err

}

// SetupComicExtra will initialize the comic based
// on comicextra.com
func SetupComicExtra(c *core.Comic) error {
	name := c.SplitURL()[3]
	issueNumber := c.SplitURL()[4]
	c.SetInfo(name, issueNumber)

	links, err := retrieveComicExtraImageLinks(c)
	c.SetImageLinks(links)

	return err
}
