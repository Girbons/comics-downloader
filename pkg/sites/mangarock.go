package sites

import (
	"strings"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/Girbons/mangarock"
	log "github.com/sirupsen/logrus"
)

type Mangarock struct {
	Client *mangarock.Client
}

func NewMangarock(options map[string]string) *Mangarock {
	return &Mangarock{
		Client: mangarock.NewClient(options),
	}
}

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

func (m *Mangarock) retrieveLastIssue(url string) (string, error) {
	url = strings.Join(util.TrimAndSplitURL(url)[:5], "/")
	series := util.TrimAndSplitURL(url)[4]

	info, infoErr := m.Client.Info(series)

	if infoErr != nil {
		return "", infoErr
	}

	lastChapter := info.Data.Chapters[len(info.Data.Chapters)-1]
	chapterUrl := url + "/chapter/" + lastChapter.OID

	return chapterUrl, nil
}

// RetrieveIssueLinks gets a slice of urls for all issues in a comic
func (m *Mangarock) RetrieveIssueLinks(url string, all, last bool) ([]string, error) {
	if last {
		url, err := m.retrieveLastIssue(url)
		return []string{url}, err
	}

	if all && m.isSingleIssue(url) {
		url = strings.Join(util.TrimAndSplitURL(url)[:5], "/")
	} else if m.isSingleIssue(url) {
		return []string{url}, nil
	}

	var links []string

	series := util.TrimAndSplitURL(url)[4]
	// get info about the manga
	info, infoErr := m.Client.Info(series)
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

func (m *Mangarock) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	series := parts[4]
	chapterID := parts[6]

	info, err := m.Client.Info(series)

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
	// get info about the manga
	parts := util.TrimAndSplitURL(comic.URLSource)
	series := parts[4]
	chapterID := parts[6]
	// retrieve pages
	info, err := m.Client.Info(series)
	if err != nil {
		return err
	}

	pages, pagesErr := m.Client.Pages(chapterID)

	comic.Author = info.Data.Author
	comic.Links = pages.Data

	return pagesErr
}
