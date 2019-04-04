package sites

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/mangarock"
	log "github.com/sirupsen/logrus"
)

func findChapterName(chapterID string, chapters []*mangarock.Chapter) (string, bool) {
	for _, chapter := range chapters {
		if chapter.OID == chapterID {
			return chapter.Name, true
		}
	}
	return "", false
}

func SetupMangaRock(c *core.Comic) error {
	series := c.SplitURL()[4]
	chapterID := c.SplitURL()[6]

	client := mangarock.NewClient()
	if _, ok := c.Options["country"]; ok {
		client.SetOptions(c.Options)
	}
	// get info about the manga
	info, infoErr := client.Info(series)
	if infoErr != nil {
		log.WithFields(log.Fields{
			"series": series,
			"source": c.Source,
		}).Error(infoErr)
	}
	// retrieve pages
	pages, pagesErr := client.Pages(chapterID)

	chapter, found := findChapterName(chapterID, info.Data.Chapters)
	if !found {
		log.WithFields(log.Fields{
			"source": c.Source,
		}).Warning("Chapter not found")
		chapter = chapterID
	}

	c.SetInfo(info.Data.Name, chapter)
	c.SetAuthor(info.Data.Author)
	c.SetImageLinks(pages.Data)

	return pagesErr
}
