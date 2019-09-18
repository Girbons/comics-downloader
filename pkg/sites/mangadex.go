package sites

import (
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/bake/mangadex"
)

type Mangadex struct {
	Client *mangadex.Client
}

// NewMangadex returns a Mangadex instance
func NewMangadex() *Mangadex {
	return &Mangadex{
		Client: mangadex.New(),
	}
}

func (m *Mangadex) getLinks(url string) ([]string, error) {
	var links []string

	parts := util.TrimAndSplitURL(url)
	chapterID := parts[4]

	res, err := m.Client.Chapter(chapterID)
	if err != nil {
		return links, err
	}

	links = res.Images()

	return links, nil
}

func (m *Mangadex) RetrieveIssueLinks(url string, all, last bool) ([]string, error) {
	return []string{url}, nil
}

func (m *Mangadex) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	chapterID := parts[4]

	chapterInfo, err := m.Client.Chapter(chapterID)
	if err != nil {
		return "", ""
	}

	mangaInfo, _, err := m.Client.Manga(string(chapterInfo.MangaID))
	if err != nil {
		return "", ""
	}

	issueNumber := fmt.Sprintf("Vol %s Chapter %s, %s", chapterInfo.Volume, chapterInfo.Chapter, chapterInfo.Title)

	return mangaInfo.Title, issueNumber
}

// Initialize loads links and metadata from mangadex
func (m *Mangadex) Initialize(comic *core.Comic) error {
	links, err := m.getLinks(comic.URLSource)
	comic.Links = links

	return err
}
