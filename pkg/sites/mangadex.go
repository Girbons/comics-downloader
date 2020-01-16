package sites

import (
	"errors"
	"fmt"

	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/bake/mangadex"
)

type Mangadex struct {
	country string
	baseURL string
	Client  *mangadex.Client
}

// NewMangadex returns a Mangadex instance
func NewMangadex(country, source string) *Mangadex {
	mangadexBase := "https://"+source+"/"
	return &Mangadex{
		country: country,
		baseURL: mangadexBase,
		Client: mangadex.New(
			mangadex.WithBase(mangadexBase),
		),
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
	parts := util.TrimAndSplitURL(url)
	if len(parts) < 5 {
		return nil, errors.New("URL not supported")
	}
	switch parts[3] {
	case "chapter":
		return []string{url}, nil
	case "title":
		_, cs, err := m.Client.Manga(parts[4])
		if err != nil {
			return nil, err
		}
		var urls []string
		for _, c := range cs {
			if m.country != "" && c.LangCode != m.country {
				continue
			}
			urls = append(urls, m.baseURL+"chapter"+c.ID.String())
		}
		if len(urls) == 0 {
			return nil, errors.New("no chapters found")
		}
		if last {
			urls = urls[len(urls)-1:]
		}
		return urls, nil
	default:
		return nil, errors.New("URL not supported")
	}
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
