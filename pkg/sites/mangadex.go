package sites

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
	"github.com/Girbons/comics-downloader/pkg/util"
)

const (
	mangadexAPIBase        = "https://api.mangadex.org"
	mangadexChapterBase    = "https://mangadex.org/chapter"
	mangadexUploadsBase    = "https://uploads.mangadex.org/data"
	mangadexRequestTimeout = 8 * time.Second
)

// Mangadex represents a mangadex instance.
type Mangadex struct {
	country     string
	options     *config.Options
	client      *httpclient.ComicClient
	apiBase     string
	chapterBase string
	uploadsBase string
}

// NewMangadex returns a Mangadex instance.
func NewMangadex(options *config.Options) *Mangadex {
	client := options.Client
	if client == nil {
		client = httpclient.NewComicClient()
		options.Client = client
	}

	return &Mangadex{
		country:     strings.ToLower(options.Country),
		options:     options,
		client:      client,
		apiBase:     mangadexAPIBase,
		chapterBase: mangadexChapterBase,
		uploadsBase: mangadexUploadsBase,
	}
}

func (m *Mangadex) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), mangadexRequestTimeout)
}

func joinURL(base, suffix string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(suffix, "/")
}

func (m *Mangadex) getManga(mangaID string) (string, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/manga/%s", mangaID))
	var mangaRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				Titles map[string]string `json:"title"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := fetchJSON(ctx, m.client, endpoint, &mangaRes); err != nil {
		return "", err
	}
	if strings.ToLower(mangaRes.Result) != "ok" {
		return "", fmt.Errorf("unexpected response")
	}

	for lang, t := range mangaRes.Data.Attributes.Titles {
		if m.country == "" || m.country == strings.ToLower(lang) {
			return t, nil
		}
	}

	// Fallback to any available title.
	for _, t := range mangaRes.Data.Attributes.Titles {
		return t, nil
	}

	return "", fmt.Errorf("no title found for manga %s", mangaID)
}

// getChapters fetches chapter URLs for the given manga.
func (m *Mangadex) getChapters(mangaID string) ([]string, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/manga/%s/aggregate", mangaID))
	if m.country != "" {
		q := url.Values{}
		q.Add("translatedLanguage[]", m.country)
		endpoint += "?" + q.Encode()
	}

	body, err := fetchBytes(ctx, m.client, endpoint)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 || body[0] == '[' {
		return []string{}, nil
	}

	var chaptersRes struct {
		Result  string `json:"result"`
		Volumes map[string]struct {
			Chapters map[string]struct {
				ID   string `json:"id"`
				Name string `json:"chapter"`
			} `json:"chapters"`
		} `json:"volumes"`
	}

	if err := json.Unmarshal(body, &chaptersRes); err != nil {
		return nil, err
	}
	if strings.ToLower(chaptersRes.Result) != "ok" {
		return nil, fmt.Errorf("unexpected response")
	}

	var ids []string
	for _, v := range chaptersRes.Volumes {
		for _, c := range v.Chapters {
			ids = append(ids, joinURL(m.chapterBase, c.ID))
		}
	}
	return ids, nil
}

// getChapter retrieves metadata and image links for a single chapter.
func (m *Mangadex) getChapter(chapterID string) (mangaID, volume, chapter, title string, images []string, err error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/chapter/%s", chapterID))
	var chapterRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				Volume  string `json:"volume"`
				Chapter string `json:"chapter"`
				Title   string `json:"title"`
			} `json:"attributes"`
			Relationships []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := fetchJSON(ctx, m.client, endpoint, &chapterRes); err != nil {
		return "", "", "", "", nil, err
	}
	if strings.ToLower(chapterRes.Result) != "ok" {
		return "", "", "", "", nil, fmt.Errorf("unexpected response")
	}

	imagesEndpoint := joinURL(m.apiBase, fmt.Sprintf("/at-home/server/%s", chapterID))
	var imagesRes struct {
		Result  string `json:"result"`
		Chapter struct {
			Hash string   `json:"hash"`
			Data []string `json:"data"`
		} `json:"chapter"`
	}

	if err := fetchJSON(ctx, m.client, imagesEndpoint, &imagesRes); err != nil {
		return "", "", "", "", nil, err
	}
	if strings.ToLower(imagesRes.Result) != "ok" {
		return "", "", "", "", nil, fmt.Errorf("unexpected response")
	}

	for _, file := range imagesRes.Chapter.Data {
		imageURL := joinURL(m.uploadsBase, fmt.Sprintf("%s/%s", imagesRes.Chapter.Hash, file))
		images = append(images, imageURL)
	}

	if m.options.Debug && len(images) > 0 && m.options.Logger != nil {
		m.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(images, " ")))
	}

	for _, rel := range chapterRes.Data.Relationships {
		if rel.Type == "manga" {
			mangaID = rel.ID
			break
		}
	}

	return mangaID, chapterRes.Data.Attributes.Volume, chapterRes.Data.Attributes.Chapter, chapterRes.Data.Attributes.Title, images, nil
}

// RetrieveIssueLinks retrieve the issue links for the given comic.
func (m *Mangadex) RetrieveIssueLinks() ([]string, error) {
	parts := util.TrimAndSplitURL(m.options.URL)
	if len(parts) < 5 {
		return nil, errors.New("URL not supported")
	}
	switch parts[3] {
	case "chapter":
		return []string{m.options.URL}, nil
	case "title":
		return m.getChapters(parts[4])
	default:
		return nil, errors.New("URL not supported")
	}
}

// GetInfo extracts the basic info from the given url.
func (m *Mangadex) GetInfo(urlValue string) (string, string) {
	parts := util.TrimAndSplitURL(urlValue)
	if len(parts) < 5 {
		return "", ""
	}
	switch parts[3] {
	case "chapter":
		mangaID, volume, chapter, title, _, err := m.getChapter(parts[4])
		if err != nil {
			return "", ""
		}
		chapterTitle := fmt.Sprintf("Vol %s Chapter %s", volume, chapter)
		if title != "" {
			chapterTitle += fmt.Sprintf(", %s", title)
		}
		mangaTitle, err := m.getManga(mangaID)
		if err != nil {
			return "", chapterTitle
		}
		return mangaTitle, chapterTitle
	case "title":
		mangaTitle, err := m.getManga(parts[4])
		if err != nil {
			return "", ""
		}
		return mangaTitle, ""
	default:
		return "", ""
	}
}

// Initialize loads links and metadata from mangadex.
func (m *Mangadex) Initialize(comic *core.Comic) error {
	parts := util.TrimAndSplitURL(comic.URLSource)
	if len(parts) < 5 {
		return fmt.Errorf("URL not supported")
	}
	_, _, _, _, images, err := m.getChapter(parts[4])
	if err != nil {
		return err
	}
	comic.Links = images
	return nil
}
