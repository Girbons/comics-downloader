package mangarock

import (
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
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

func RetrieveIssueLinks(url string, options map[string]string) ([]string, error) {
	if isSingleIssue(url) {
		return []string{url}, nil
	}

	var links []string

	series := util.SplitURL(url)[4]

	client := mangarock.NewClient()
	if _, ok := options["country"]; ok {
		client.SetOptions(options)
	}
	// get info about the manga
	info, infoErr := client.Info(series)
	if infoErr != nil {
		log.Error(infoErr)
	}

	for _, chapter := range info.Data.Chapters {
		chapterUrl := url + "/chapter/" + chapter.OID
		if util.IsURLValid(chapterUrl) {
			links = append(links, chapterUrl)
		}
	}

	return links, nil
}

func isSingleIssue(url string) bool {
	return len(util.SplitURL(url)) >= 6
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
