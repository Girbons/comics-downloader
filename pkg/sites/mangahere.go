package sites

import (
	"fmt"
	"regexp"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/anaskhan96/soup"
	log "github.com/sirupsen/logrus"
)

func retrieveMangaHereImageLinks(c *core.Comic) ([]string, error) {
	var (
		pageLinks []string
		imgLinks  []string
	)
	// compile the image regex
	re := regexp.MustCompile(c.ImageRegex)

	response, err := soup.Get(c.URLSource)

	if err != nil {
		log.Error(err)
	}

	document := soup.HTMLParse(response)
	links := document.FindAll("option")

	for _, link := range links {
		newLink := fmt.Sprintf("http://%s", link.Attrs()["value"][2:])
		if !util.IsValueInSlice(newLink, pageLinks) {
			pageLinks = append(pageLinks, newLink)
		}
	}

	for _, link := range pageLinks {
		if link != "" {
			imgResponse, imgResponseError := soup.Get(link)

			if imgResponseError != nil {
				log.Error(imgResponseError)
			}

			match := re.FindAllStringSubmatch(imgResponse, -1)
			for _, f := range match {
				if util.IsUrlValid(f[1]) {
					imgLinks = append(imgLinks, f[1])
				}
			}
		}
	}
	return imgLinks, err
}

// SetupMangaHere will initialize the comic based
// on mangahere.cc
func SetupMangaHere(c *core.Comic, splittedUrl []string) error {
	name := splittedUrl[4]
	issueNumber := splittedUrl[5]
	imageRegex := `<img[^>]+src="([^">]+)"`
	c.SetInfo(name, issueNumber, imageRegex)

	links, err := retrieveMangaHereImageLinks(c)
	c.SetImageLinks(links)

	return err
}
