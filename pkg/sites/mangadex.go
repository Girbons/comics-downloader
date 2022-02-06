package sites

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/Girbons/comics-downloader/pkg/util"
)

// Mangadex represents a mangadex instance.
type Mangadex struct {
	country string
	baseURL string
	options *config.Options
}

// NewMangadex returns a Mangadex instance
func NewMangadex(options *config.Options) *Mangadex {
	return &Mangadex{
		country: strings.ToLower(options.Country),
		options: options,
	}
}

func (m *Mangadex) getManga(mangaID string) (title string, err error) {
	url := fmt.Sprintf("https://api.mangadex.org/manga/%s", mangaID)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var mangaRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				Titles map[string]string `json:"title"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&mangaRes); err != nil {
		return "", err
	}
	if mangaRes.Result != "ok" {
		return "", fmt.Errorf("Unexpected response")
	}
	for lang, t := range mangaRes.Data.Attributes.Titles {
		title = t
		if m.country == "" || m.country == lang {
			break
		}
	}
	return title, nil
}

// Get a list of chapter IDs of a manga.
func (m *Mangadex) getChapters(mangaID string) ([]string, error) {
	url := fmt.Sprintf("https://api.mangadex.org/manga/%s/aggregate?", mangaID)
	if m.country != "" {
		url += fmt.Sprintf("&translatedLanguage[]=%s", m.country)
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var chaptersRes struct {
		Result  string `json:"result"`
		Volumes map[string]struct {
			Name     string `json:"volume"`
			Chapters map[string]struct {
				ID   string `json:"id"`
				Name string `json:"chapter"`
			} `json:"chapters"`
		} `json:"volumes"`
	}
	if err := json.NewDecoder(res.Body).Decode(&chaptersRes); err != nil {
		// This is not ideal. MangaDex returns an empty array (`[]`) when no volume
		// exists and a `map[string]interface{}` otherwise.
		return []string{}, nil
	}
	if chaptersRes.Result != "ok" {
		return nil, fmt.Errorf("Unexpected response")
	}
	var ids []string
	for _, v := range chaptersRes.Volumes {
		for _, c := range v.Chapters {
			url := fmt.Sprintf("https://mangadex.org/chapter/%s", c.ID)
			ids = append(ids, url)
		}
	}
	return ids, nil
}

// Get a list of a chapters images.
func (m *Mangadex) getChapter(chapterID string) (mangaID, volume, chapter, title string, images []string, err error) {
	url := fmt.Sprintf("https://api.mangadex.org/chapter/%s", chapterID)
	res, err := http.Get(url)
	if err != nil {
		return "", "", "", "", nil, err
	}
	defer res.Body.Close()
	var chapterRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				Volume  string   `json:"volume"`
				Chapter string   `json:"chapter"`
				Title   string   `json:"title"`
				Hash    string   `json:"hash"`
				Data    []string `json:"data"`
			} `json:"attributes"`
			Relationships []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := json.NewDecoder(res.Body).Decode(&chapterRes); err != nil {
		return "", "", "", "", nil, err
	}
	if chapterRes.Result != "ok" {
		return "", "", "", "", nil, fmt.Errorf("Unexpected response")
	}

	var imagesRes struct {
		Result  string `json:"result"`
		Chapter struct {
			Hash string   `json:"hash"`
			Data []string `json:"data"`
		} `json:"chapter`
	}

	res, err = http.Get(fmt.Sprintf("https://api.mangadex.org/at-home/server/%s", chapterID))

	if err := json.NewDecoder(res.Body).Decode(&imagesRes); err != nil {
		return "", "", "", "", nil, err
	}

	for _, file := range imagesRes.Chapter.Data {
		imageUrl := fmt.Sprintf("https://uploads.mangadex.org/data/%s/%s", imagesRes.Chapter.Hash, file)
		images = append(images, imageUrl)
	}

	if m.options.Debug {
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
func (m *Mangadex) GetInfo(url string) (string, string) {
	parts := util.TrimAndSplitURL(url)
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

// Initialize loads links and metadata from mangadex
func (m *Mangadex) Initialize(comic *core.Comic) error {
	parts := util.TrimAndSplitURL(comic.URLSource)
	if len(parts) < 4 {
		return fmt.Errorf("URL not supported")
	}
	_, _, _, _, images, err := m.getChapter(parts[4])
	if err != nil {
		return err
	}
	comic.Links = images
	return nil
}
