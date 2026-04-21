package sites

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
	"github.com/Girbons/comics-downloader/pkg/util"
)

const (
	mangadexWebBase       = "https://mangadex.org"
	mangadexAPIBase       = "https://api.mangadex.org"
	mangadexUploadsBase   = "https://uploads.mangadex.org"
	mangadexChapterBase   = mangadexWebBase + "/chapter"
	mangadexUploadsData   = mangadexUploadsBase + "/data"
	mangadexUploadsCovers = mangadexUploadsBase + "/covers"
)

// Mangadex represents a mangadex instance.
type Mangadex struct {
	country       string
	options       *config.Options
	client        *httpclient.ComicClient
	apiBase       string
	chapterBase   string
	uploadsData   string
	uploadsCovers string
}

// NewMangadex returns a Mangadex instance.
func NewMangadex(options *config.Options) *Mangadex {
	return &Mangadex{
		country:       strings.ToLower(options.Country),
		options:       options,
		client:        options.Client,
		apiBase:       mangadexAPIBase,
		chapterBase:   mangadexChapterBase,
		uploadsData:   mangadexUploadsData,
		uploadsCovers: mangadexUploadsCovers,
	}
}

func init() {
	SupportedSites["mangadex"] = SupportedSite{
		IsEnabled: true,
		Loader:    func(opts *config.Options) BaseSite { return NewMangadex(opts) },
	}
}

func (m *Mangadex) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), m.options.RequestTimeout)
}

func joinURL(base, suffix string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(suffix, "/")
}

type mangadexSeries struct {
	ID string

	Title          string // the title we want to use for the series, which may be localized based on the country option
	LocalizedTitle map[string]string
	Description    map[string]string

	IsManga          bool
	OriginalLanguage string
	Year             *int
	CoverURL         *string
	Rating           *float64

	ContentRating core.AgeRating
	Tags          []string
	Genres        []string
	WebLinks      []string

	Authors []string
	Artists []string

	IsOneShot bool
}

func (m *Mangadex) getMangaCoverURL(mangaID, coverID string) (string, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/cover/%s", coverID))
	var coverRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				FileName string `json:"fileName"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := m.client.FetchJSON(ctx, endpoint, &coverRes); err != nil {
		return "", err
	}
	if strings.ToLower(coverRes.Result) != "ok" {
		return "", fmt.Errorf("unexpected response")
	}

	return joinURL(m.uploadsCovers, fmt.Sprintf("%s/%s", mangaID, coverRes.Data.Attributes.FileName)), nil
}

func (m *Mangadex) getAuthorInfo(authorID string) (string, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/author/%s", authorID))
	var authorRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := m.client.FetchJSON(ctx, endpoint, &authorRes); err != nil {
		return "", err
	}
	if strings.ToLower(authorRes.Result) != "ok" {
		return "", fmt.Errorf("unexpected response")
	}

	return authorRes.Data.Attributes.Name, nil
}

func (m *Mangadex) getMangaRating(mangaID string) (float64, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	// need to do weird [] syntax so it counts as an array to the mangadex API
	endpoint := joinURL(m.apiBase, fmt.Sprintf("/statistics/manga?manga[]=%s", mangaID))
	var ratingRes struct {
		Result     string `json:"result"`
		Statistics map[string]struct {
			Rating struct {
				Bayesian float64 `json:"bayesian"`
			} `json:"rating"`
		} `json:"statistics"`
	}
	if err := m.client.FetchJSON(ctx, endpoint, &ratingRes); err != nil {
		return 0, err
	}
	if strings.ToLower(ratingRes.Result) != "ok" {
		return 0, fmt.Errorf("unexpected response")
	}

	stats, exists := ratingRes.Statistics[mangaID]
	if !exists {
		return 0, fmt.Errorf("no statistics found for manga ID %s", mangaID)
	}

	return stats.Rating.Bayesian, nil

}

func (m *Mangadex) getMangaInfo(mangaID string) (mangadexSeries, error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	endpoint := joinURL(m.apiBase, fmt.Sprintf("/manga/%s", mangaID))
	var mangaRes struct {
		Result string `json:"result"`
		Data   struct {
			Attributes struct {
				Title                  map[string]string   `json:"title"`
				AltTitles              []map[string]string `json:"altTitles"`
				Description            map[string]string   `json:"description"`
				Links                  map[string]string   `json:"links"`
				OriginalLanguage       string              `json:"originalLanguage"`
				Year                   *int                `json:"year"`
				ContentRating          string              `json:"contentRating"`
				PublicationDemographic *string             `json:"publicationDemographic"`
				Tags                   []struct {
					ID         string `json:"id"`
					Type       string `json:"type"` // should be "tag"
					Attributes struct {
						Name        map[string]string `json:"name"`
						Description map[string]string `json:"description"`
						Group       string            `json:"group"`
						Version     int               `json:"version"`
					} `json:"attributes"`
				} `json:"tags"`
				AvailableTranslatedLanguages []string `json:"availableTranslatedLanguages"`
			} `json:"attributes"`
			Relationships []struct {
				ID   string `json:"id"`
				Type string `json:"type"` // e.g. "author", "artist", "cover_art"
			} `json:"relationships"`
		} `json:"data"`
	}

	if err := m.client.FetchJSON(ctx, endpoint, &mangaRes); err != nil {
		return mangadexSeries{}, err
	}
	if strings.ToLower(mangaRes.Result) != "ok" {
		return mangadexSeries{}, fmt.Errorf("unexpected response")
	}

	manga := mangadexSeries{
		LocalizedTitle: mangaRes.Data.Attributes.Title,
		Description:    mangaRes.Data.Attributes.Description,
		IsManga:        true, // default to true since it's a manga site
	}

	// Set titles
	foundTitle := false
	for lang, title := range mangaRes.Data.Attributes.Title {
		if m.country == "" || m.country == strings.ToLower(lang) {
			manga.Title = title
			foundTitle = true
			break
		}

		// manga.LocalizedTitle[lang] = title
	}

	// try and fill in any missing localized titles
	for _, grouping := range mangaRes.Data.Attributes.AltTitles {
		for lang, title := range grouping {
			if _, exists := manga.LocalizedTitle[lang]; !exists {
				manga.LocalizedTitle[lang] = title
			}

			// if main title is missing, try to use the alt titles to fill it
			if !foundTitle && (m.country == "" || m.country == strings.ToLower(lang)) {
				manga.Title = title
				foundTitle = true
			}
		}
	}

	if !foundTitle {
		// If still haven't found anything fallback to any available title.
		for _, title := range mangaRes.Data.Attributes.Title {
			manga.Title = title
		}
	}

	// TODO: should we error if the manga isn't available in the specified country?
	// previously we have just let the process continue
	// if m.country != "" && !slices.Contains(mangaRes.Data.Attributes.AvailableTranslatedLanguages, strings.ToLower(m.country)) {
	// 	return mangadexSeries{}, fmt.Errorf("manga \"%s\" not available in country %s", manga.Title, m.country)
	// }

	// get author, artist, and cover info from relationships
	cachedAuthors := map[string]string{} // cache author info to avoid duplicate requests
	for _, rel := range mangaRes.Data.Relationships {
		switch rel.Type {
		case "cover_art":
			coverURL, err := m.getMangaCoverURL(mangaID, rel.ID)
			if err != nil {
				return mangadexSeries{}, err
			}
			manga.CoverURL = &coverURL
		case "author":
			var err error
			authorName, ok := cachedAuthors[rel.ID]
			if !ok {
				authorName, err = m.getAuthorInfo(rel.ID)
				if err != nil {
					return mangadexSeries{}, err
				}
			}
			manga.Authors = append(manga.Authors, authorName)
			cachedAuthors[rel.ID] = authorName
		case "artist":
			var err error
			artistName, ok := cachedAuthors[rel.ID]
			if !ok {
				artistName, err = m.getAuthorInfo(rel.ID)
				if err != nil {
					return mangadexSeries{}, err
				}
			}
			manga.Artists = append(manga.Artists, artistName)
			cachedAuthors[rel.ID] = artistName
		}
	}

	// get rating info
	rating, err := m.getMangaRating(mangaID)
	if err != nil {
		return mangadexSeries{}, err
	}
	manga.Rating = &rating

	// handle tags and genres
	if mangaRes.Data.Attributes.PublicationDemographic != nil {
		manga.Tags = append(manga.Tags, *mangaRes.Data.Attributes.PublicationDemographic)
	}
	for _, tag := range mangaRes.Data.Attributes.Tags {
		// using English name for tags since they usally don't have localized names
		name, ok := tag.Attributes.Name["en"]
		if !ok {
			// if no English name, try to use the first available name
			for _, n := range tag.Attributes.Name {
				name = n
				break
			}
		}

		if name == "" {
			continue
		}

		// classify tag groups
		switch tag.Attributes.Group {
		case "genre":
			manga.Genres = append(manga.Genres, name)
		case "tag", "content", "theme":
			manga.Tags = append(manga.Tags, name)
		case "format":
			if name == "Oneshot" {
				manga.IsOneShot = true
			}
			manga.Tags = append(manga.Tags, name)
		}
	}

	// see https://api.mangadex.org/docs/3-enumerations/#manga-content-rating
	switch mangaRes.Data.Attributes.ContentRating {
	case "safe":
		manga.ContentRating = core.AgeRatingEveryone
	case "suggestive", "erotica":
		manga.ContentRating = core.AgeRatingMature
	case "pornographic":
		manga.ContentRating = core.AgeRatingAO18
	default:
		manga.ContentRating = core.AgeRatingEveryone
	}

	// set link to main mangadex page for the manga
	manga.WebLinks = append(manga.WebLinks, fmt.Sprintf("%s/title/%s", mangadexWebBase, mangaID))
	// also add any additional links provided by mangadex
	for key, link := range mangaRes.Data.Attributes.Links {
		fullURL := m.mangaLinkToFullURL(key, link)
		if fullURL != "" {
			manga.WebLinks = append(manga.WebLinks, fullURL)
		}
	}
	manga.OriginalLanguage = mangaRes.Data.Attributes.OriginalLanguage
	manga.Year = mangaRes.Data.Attributes.Year

	return manga, nil
}

func (m *Mangadex) mangaLinkToFullURL(key, link string) string {
	switch key {
	case "al":
		return "https://www.anilist.co/manga/" + link
	case "ap":
		return "https://www.animeplanet.com/manga/" + link
	case "bw":
		return "https://www.bookwalker.jp/" + link
	case "mu":
		return "https://www.mangaupdates.com/series.html?id=" + link
	case "nu":
		return "https://www.novelupdates.com/series/" + link
	case "kt":
		// if int use id
		if _, err := strconv.Atoi(link); err == nil {
			return "https://kitsu.io/api/edge/manga/" + link
		}
		// else use slug
		return "https://kitsu.io/api/edge/manga?filter[slug]=" + link
	case "amz":
		return link
	case "ebj":
		return link
	case "mal":
		return "https://myanimelist.net/manga/" + link
	case "cdj":
		return link
	case "raw":
		return link
	case "engtl":
		return link
	default:
		return ""
	}
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

	body, err := m.client.FetchBytes(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 || body[0] == '[' {
		return []string{}, nil
	}

	var chaptersRes struct {
		Result  string `json:"result"`
		Volumes map[string]struct {
			Volume   string `json:"volume"` // volume name, or "none"
			Count    int    `json:"count"`
			Chapters map[string]struct {
				Chapter       string   `json:"chapter"` // the chapter number, not the chapter name
				ID            string   `json:"id"`
				IsUnavailable bool     `json:"isUnavailable"`
				Others        []string `json:"others"` // list of alternative chapter IDs for the same chapter (e.g. for different languages)
				// TODO: add check to try and find target language in others if the main chapter is not what we want
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
			if c.IsUnavailable {
				continue
			}

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

// getChapterInfo retrieves metadata and image links for a single chapter.
func (m *Mangadex) getChapterInfo(chapterID string) (chapterInfo mangadexChapter, err error) {
	ctx, cancel := m.requestContext()
	defer cancel()

	// TODO: set total number of chapters available in the manga
	// ~~TODO: set volume number in metadata if available~~ Seems to be a non-issue??
	// ~~TODO: set issue number to chapterNum~~ again seesm to be a non-issue?

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

	if err := m.client.FetchJSON(ctx, endpoint, &chapterRes); err != nil {
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

	if err := m.client.FetchJSON(ctx, imagesEndpoint, &imagesRes); err != nil {
		return mangadexChapter{}, err
	}
	if strings.ToLower(imagesRes.Result) != "ok" {
		return mangadexChapter{}, fmt.Errorf("unexpected response")
	}

	var imageLinks []string
	for _, file := range imagesRes.Chapter.Data {
		imageURL := joinURL(m.uploadsData, fmt.Sprintf("%s/%s", imagesRes.Chapter.Hash, file))
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
func (m *Mangadex) GetInfo(urlValue string) (string, string, error) {
	parts := util.TrimAndSplitURL(urlValue)
	if len(parts) < 5 {
		return "", "", errors.New("URL not supported")
	}
	switch parts[3] {
	case "chapter":
		chapter, err := m.getChapterInfo(parts[4])
		if err != nil {
			return "", "", err
		}

		manga, err := m.getMangaInfo(chapter.MangaID)
		if err != nil {
			return "", "", err
		}
		return manga.Title, chapter.ChapterNumber, nil

	case "title":
		manga, err := m.getMangaInfo(parts[4])
		if err != nil {
			return "", "", err
		}
		return manga.Title, "", nil
	default:
		return "", "", errors.New("URL not supported")
	}
}

// Initialize loads links and metadata from mangadex.
func (m *Mangadex) Initialize(comic *core.ComicIssue) error {

	// TODO: fix localized title being set as "Series Title" in metadata
	// TODO: fix chapter title being set to a fs compatible sanatized version of the chapter name instead of the actual chapter name

	if comic == nil {
		return fmt.Errorf("comic is nil")
	}
	if comic.Source == nil {
		return fmt.Errorf("comic source is nil")
	}

	parts := util.TrimAndSplitURL(comic.Source.URL)
	if len(parts) < 5 {
		return fmt.Errorf("URL not supported")
	}
	chapter, err := m.getChapterInfo(parts[4])
	if err != nil {
		return err
	}
	manga, err := m.getMangaInfo(chapter.MangaID)
	if err != nil {
		return err
	}

	comic.ChapterName = chapter.ChapterTitle // changing the title seems to break path resolving for some reason, probably because the folder has already been created by the time we get to this point, so we just keep the name as is until the metadata system is reworked
	comic.IssueNumber = chapter.ChapterNumber
	comic.Volume = chapter.Volume
	comic.LanguageISO = &chapter.TranslatedLanguage
	comic.ReleaseDate = &chapter.PublishAt

	// ensure metadata object exists
	if comic.SeriesMetadata == nil {
		comic.SeriesMetadata = &core.SeriesMetadata{}
	}

	comic.SeriesMetadata.Title = manga.Title
	comic.SeriesMetadata.LocalizedTitle = manga.LocalizedTitle
	comic.SeriesMetadata.Description = manga.Description
	comic.SeriesMetadata.Genres = manga.Genres
	comic.SeriesMetadata.Tags = manga.Tags
	comic.SeriesMetadata.AgeRating = &manga.ContentRating
	comic.SeriesMetadata.WebLinks = manga.WebLinks
	comic.SeriesMetadata.IsManga = &manga.IsManga
	comic.SeriesMetadata.CommunityRating = manga.Rating

	if manga.IsOneShot {
		format := core.ComicFormatOneShot
		comic.ComicFormat = &format
	}

	for _, author := range manga.Authors {
		comic.SeriesMetadata.Creators = append(comic.SeriesMetadata.Creators, core.SeriesCreator{
			Name: author,
			Role: core.CreatorRoleWriter,
		})
	}
	for _, artist := range manga.Artists {
		comic.SeriesMetadata.Creators = append(comic.SeriesMetadata.Creators, core.SeriesCreator{
			Name: artist,
			Role: core.CreatorRolePenciller,
		})
	}

	comic.SeriesMetadata.CoverURL = manga.CoverURL
	comic.ImageLinks = chapter.ImageLinks

	return nil
}
