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

type mangadexChapter struct {
	ChapterID     string
	ChapterNumber string
	ChapterTitle  string
	Volume        *string // can be null if not available

	TranslatedLanguage string // Not sure if this always exists or not e.g. "en", "jp"
	PublishAt          time.Time

	MangaID string

	ImageLinks []string
}

// getChapter retrieves metadata and image links for a single chapter.
func (m *Mangadex) getChapter(chapterID string) (chapterInfo mangadexChapter, err error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/chapter/%s", chapterID))
	var chapterRes struct {
		Result string `json:"result"`
		Data   struct {
			ID         string `json:"id"`   // chapter ID
			Type       string `json:"type"` // should be "chapter"
			Attributes struct {
				Volume             *string `json:"volume"` // is null if not available
				Chapter            *string `json:"chapter"`
				Title              *string `json:"title"`
				TranslatedLanguage string  `json:"translatedLanguage"`
				PublishAt          string  `json:"publishAt"`
			} `json:"attributes"`
			Relationships []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := fetchJSON(ctx, m.client, endpoint, &chapterRes); err != nil {
		return mangadexChapter{}, err
	}
	if strings.ToLower(chapterRes.Result) != "ok" {
		return mangadexChapter{}, fmt.Errorf("unexpected response")
	}

	publishedAt, err := time.Parse(time.RFC3339, chapterRes.Data.Attributes.PublishAt)
	if err != nil {
		return mangadexChapter{}, err
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
		return mangadexChapter{}, err
	}
	if strings.ToLower(imagesRes.Result) != "ok" {
		return mangadexChapter{}, fmt.Errorf("unexpected response")
	}

	var imageLinks []string
	for _, file := range imagesRes.Chapter.Data {
		imageURL := joinURL(m.uploadsBase, fmt.Sprintf("%s/%s", imagesRes.Chapter.Hash, file))
		imageLinks = append(imageLinks, imageURL)
	}

	if m.options.Debug && len(imageLinks) > 0 && m.options.Logger != nil {
		m.options.Logger.Debug(fmt.Sprintf("Image Links found: %s", strings.Join(imageLinks, " ")))
	}

	var mangaID string
	for _, rel := range chapterRes.Data.Relationships {
		if rel.Type == "manga" {
			mangaID = rel.ID
			break
		}
	}

	// just default them to empty string if not found
	var chapterNumber string
	var chapterTitle string
	if chapterRes.Data.Attributes.Chapter != nil {
		// TODO: consider defaulting to "oneshot"?
		chapterNumber = *chapterRes.Data.Attributes.Chapter
	}
	if chapterRes.Data.Attributes.Title != nil {
		chapterTitle = *chapterRes.Data.Attributes.Title
	}

	return mangadexChapter{
		ChapterID:     chapterRes.Data.ID,
		ChapterNumber: chapterNumber,
		ChapterTitle:  chapterTitle,
		Volume:        chapterRes.Data.Attributes.Volume,

		TranslatedLanguage: chapterRes.Data.Attributes.TranslatedLanguage,
		PublishAt:          publishedAt,

		MangaID: mangaID,

		ImageLinks: imageLinks,
	}, nil
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
		chapter, err := m.getChapter(parts[4])
		if err != nil {
			return "", ""
		}

		var chapterTitle string
		if chapter.Volume != nil {
			volume := *chapter.Volume
			chapterTitle = fmt.Sprintf("Vol %s Chapter %s", volume, chapter.ChapterNumber)
		} else {
			chapterTitle = fmt.Sprintf("Chapter %s", chapter.ChapterNumber)
		}

		if chapter.ChapterTitle != "" {
			chapterTitle += fmt.Sprintf(", %s", chapter.ChapterTitle)
		}
		mangaTitle, err := m.getManga(chapter.MangaID)
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
func (m *Mangadex) Initialize(comic *core.ComicIssue) error {
	parts := util.TrimAndSplitURL(comic.Source.URL)
	if len(parts) < 5 {
		return fmt.Errorf("URL not supported")
	}
	chapter, err := m.getChapter(parts[4])
	if err != nil {
		return err
	}

	// comic.Name = chapter.ChapterTitle // changing the title seems to break path resolving for some reason, probably because the folder has already been created by the time we get to this point, so we just keep the name as is until the metadata system is reworked
	comic.IssueNumber = chapter.ChapterNumber
	comic.Volume = chapter.Volume
	comic.LanguageISO = &chapter.TranslatedLanguage
	comic.ReleaseDate = &chapter.PublishAt

	comic.ImageLinks = chapter.ImageLinks

	return nil
}
