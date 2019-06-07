package mangarock

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/Girbons/mangarock"
	log "github.com/sirupsen/logrus"
)

type Mangarock struct{}

func (m *Mangarock) findChapterName(chapterID string, chapters []*mangarock.Chapter) (string, bool) {
	for _, chapter := range chapters {
		if chapter.OID == chapterID {
			return chapter.Name, true
		}
	}
	return "", false
}

func (m *Mangarock) isSingleIssue(url string) bool {
	return len(util.TrimAndSplitURL(url)) >= 6
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (m *Mangarock) RetrieveIssueLinks(url string, all bool, options map[string]string) ([]string, error) {
	if all && m.isSingleIssue(url) {
		url = strings.Join(util.TrimAndSplitURL(url)[:5], "/")
	} else if m.isSingleIssue(url) {
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

func (m *Mangarock) GetInfo(url string, options map[string]string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	series := parts[4]
	chapterID := parts[6]

	client := mangarock.NewClient()

	if _, ok := options["country"]; ok {
		client.SetOptions(options)
	}

	info, err := client.Info(series)

	if err != nil {
		return "", ""
	}

	chapter, found := m.findChapterName(chapterID, info.Data.Chapters)
	if !found {
		log.Warning("Chapter not found")
		chapter = chapterID
	}

	return info.Data.Name, chapter
}

// Initialize loads links and metadata from mangarock
func (m *Mangarock) Initialize(comic *core.Comic) error {
	client := mangarock.NewClient()
	if _, ok := comic.Options["country"]; ok {
		client.SetOptions(comic.Options)
	}
	// get info about the manga
	parts := util.TrimAndSplitURL(comic.URLSource)
	series := parts[4]
	chapterID := parts[6]
	// retrieve pages
	info, err := client.Info(series)
	if err != nil {
		return err
	}

	pages, pagesErr := client.Pages(chapterID)

	comic.Author = info.Data.Author
	comic.Links = pages.Data

	return pagesErr
}
