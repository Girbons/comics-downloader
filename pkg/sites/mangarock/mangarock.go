package mangarock

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

// Initialize loads links and metadata from mangarock
func Initialize(comic *core.Comic) error {
	series := comic.SplitURL()[4]
	chapterID := comic.SplitURL()[6]

	client := mangarock.NewClient()
	if _, ok := comic.Options["country"]; ok {
		client.SetOptions(comic.Options)
	}
	// get info about the manga
	info, infoErr := client.Info(series)
	if infoErr != nil {
		log.Error(infoErr)
	}
	// retrieve pages
	pages, pagesErr := client.Pages(chapterID)

	chapter, found := findChapterName(chapterID, info.Data.Chapters)
	if !found {
		log.Warning("Chapter not found")
		chapter = chapterID
	}

	comic.Name = info.Data.Name
	comic.IssueNumber = chapter
	comic.Author = info.Data.Author
	comic.Links = pages.Data

	return pagesErr
}
