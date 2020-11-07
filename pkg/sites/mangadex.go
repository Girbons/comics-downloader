package sites

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/bake/mangadex"
)

type Mangadex struct {
	country string
	baseURL string
	Client  *mangadex.Client
	options *config.Options
}

// NewMangadex returns a Mangadex instance
func NewMangadex(options *config.Options) *Mangadex {
	mangadexBase := "https://" + options.Source + "/"
	return &Mangadex{
		country: strings.ToLower(options.Country),
		baseURL: mangadexBase,
		options: options,
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

	if m.options.Debug {
		m.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(links, " ")))
	}

	return links, nil
}

func (m *Mangadex) RetrieveIssueLinks() ([]string, error) {
	parts := util.TrimAndSplitURL(m.options.Url)
	if len(parts) < 5 {
		return nil, errors.New("URL not supported")
	}
	switch parts[3] {
	case "chapter":
		return []string{m.options.Url}, nil
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
			urls = append(urls, m.baseURL+"chapter/"+c.ID.String())
		}
		if len(urls) == 0 {
			return nil, errors.New("no chapters found")
		}
		if m.options.Last {
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

	mangaID := chapterInfo.MangaID.Number
	mangaInfo, _, err := m.Client.Manga(string(mangaID))
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
