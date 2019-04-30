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

func isSingleIssue(url string) bool {
	return len(util.TrimAndSplitURL(url)) >= 6
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func RetrieveIssueLinks(url string, options map[string]string) ([]string, error) {
	if isSingleIssue(url) {
		return []string{url}, nil
	}

	var links []string

	series := util.TrimAndSplitURL(url)[4]

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

// Initialize loads links and metadata from mangarock
func Initialize(comic *core.Comic) error {
	parts := util.TrimAndSplitURL(comic.URLSource)
	series := parts[4]
	chapterID := parts[6]

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
