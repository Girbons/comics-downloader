package sites

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
	"github.com/bake/mangadex/v2"
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
	chapterID, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, err
	}

	res, err := m.Client.Chapter(context.Background(), chapterID, nil)
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
	mangaID, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, err
	}
	switch parts[3] {
	case "chapter":
		return []string{m.options.Url}, nil
	case "title":
		cs, err := m.Client.MangaChapters(context.Background(), mangaID, nil)
		if err != nil {
			return nil, err
		}
		var urls []string
		for _, c := range cs {
			if m.country != "" && c.Language != m.country {
				continue
			}
			urls = append(urls, fmt.Sprintf("%schapter/%d", m.baseURL, c.ID))
		}
		if len(urls) == 0 {
			return nil, errors.New("no chapters found")
		}
		if m.options.Last {
			urls = urls[:1]
		}
		return urls, nil
	default:
		return nil, errors.New("URL not supported")
	}
}

func (m *Mangadex) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
	chapterID, err := strconv.Atoi(parts[4])
	if err != nil {
		return "", ""
	}

	chapterInfo, err := m.Client.Chapter(context.Background(), chapterID, nil)
	if err != nil {
		return "", ""
	}

	mangaID := chapterInfo.MangaID
	mangaInfo, err := m.Client.Manga(context.Background(), mangaID, nil)
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
