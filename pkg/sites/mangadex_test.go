package sites

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
	"github.com/stretchr/testify/require"
)

func setupMangadexServer() *httptest.Server {
	hasMangaParam := func(r *http.Request, series string) bool {
		for _, id := range r.URL.Query()["manga[]"] {
			if id == series {
				return true
			}
		}
		return false
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/statistics/manga" && hasMangaParam(r, "series-1"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"result":"ok","statistics":{"series-1":{"rating":{"bayesian":4.5}}}}`)
		case strings.HasPrefix(r.URL.Path, "/manga/series-1/aggregate"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"volumes":{
					"1":{
						"chapters":{
							"1":{"id":"chapter-1","chapter":"1"}
						}
					}
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/manga/series-1"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"data":{
					"attributes":{
						"title":{"jp":"テスト"},
						"altTitles":[
							{"en":"Test Manga"},
							{"fr":"Manga de Test"},
							{"zh": "测试漫画"}
						],
						"description":{"en":"Test manga description"},
						"publicationDemographic":"shounen",
						"contentRating":"safe",
						"tags":[
							{
								"id":"tag-1",
								"type":"tag",
								"attributes":{
									"name":{"en":"Romance"},
									"group":"genre",
									"version":1
								}
							},
							{
								"id":"tag-2",
								"type":"tag",
								"attributes":{
									"name":{"en":"School Life"},
									"group":"theme",
									"version":1
								}
							},
							{
								"id":"tag-3",
								"type":"tag",
								"attributes":{
									"name":{"en":"Doujinshi"},
									"group":"format",
									"version":1
								}
							}
						],
						"links":{
							"al": "30642"
						}
					}
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/chapter/chapter-1"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"data":{
					"attributes":{
						"volume":"1",
						"chapter":"1",
						"title":"Start",
						"publishAt":"2026-03-06T14:03:52.000Z",
						"translatedLanguage":"en"
					},
					"relationships":[{"id":"series-1","type":"manga"}]
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/at-home/server/chapter-1"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"chapter":{"hash":"HASH","data":["001.png","002.png"]}
			}`)
		case r.URL.Path == "/statistics/manga" && hasMangaParam(r, "series-2"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"result":"ok","statistics":{"series-1":{"rating":{"bayesian":4.5}}}}`)
		case strings.HasPrefix(r.URL.Path, "/manga/series-2/aggregate"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"volumes":{
					"1":{
						"chapters":{
							"1":{"id":"chapter-1","chapter":"1"}
						}
					}
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/manga/series-2"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"data":{
					"attributes":{
						"title":{"jp":"テスト"},
						"altTitles":[
							{"en":"Test Manga"}
						],
					}
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/chapter/chapter-2"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"data":{
					"attributes":{
						"volume":"1",
						"chapter":"1",
						"title":"Start",
						"publishAt":"2026-03-06T14:03:52.000Z",
						"translatedLanguage":"en"
					},
					"relationships":[{"id":"series-2","type":"manga"}]
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/at-home/server/chapter-2"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"chapter":{"hash":"HASH","data":["001.png","002.png"]}
			}`)
		default:
			http.NotFound(w, r)
		}
	}))
}

func newTestMangadex(t *testing.T, series, lang string) (*Mangadex, func()) {
	t.Helper()

	server := setupMangadexServer()

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
	)

	opts := &config.Options{
		URL:            server.URL + fmt.Sprintf("/title/%s/naruto", series),
		Country:        lang,
		SourceName:     "mangadex.org",
		Logger:         logger.NewLogger(false, nil),
		Client:         client,
		RequestTimeout: config.DefaulltRequestTimeout,
	}

	md := NewMangadex(opts)
	md.apiBase = server.URL
	md.chapterBase = server.URL + "/chapter"
	md.uploadsData = server.URL + "/data"

	cleanup := func() {
		server.Close()
	}

	return md, cleanup
}

func TestMangadexRetrieveIssueLinks(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "en")
	defer cleanup()

	md.options.All = true

	links, err := md.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{md.chapterBase + "/chapter-1"}, links)
}

func TestMangadexInitialize(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "en")
	defer cleanup()

	comic := &core.ComicIssue{
		Source: &core.ComicSource{Name: "test-source", URL: md.chapterBase + "/chapter-1"},
	}
	err := md.Initialize(comic)
	require.NoError(t, err)
	require.Equal(t, []string{
		md.uploadsData + "/HASH/001.png",
		md.uploadsData + "/HASH/002.png",
	}, comic.ImageLinks)
}

func TestMangadexInitializeMinimum(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-2", "en")
	defer cleanup()

	comic := &core.ComicIssue{
		Source: &core.ComicSource{Name: "test-source", URL: md.chapterBase + "/chapter-1"},
	}
	err := md.Initialize(comic)
	require.NoError(t, err)
	require.Equal(t, []string{
		md.uploadsData + "/HASH/001.png",
		md.uploadsData + "/HASH/002.png",
	}, comic.ImageLinks)
}

func TestMangadexGetInfoEnglish(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "en")
	defer cleanup()

	title, chapter, err := md.GetInfo(md.chapterBase + "/chapter-1")
	require.NoError(t, err)
	require.Equal(t, "Test Manga", title)
	require.Equal(t, "1", chapter)
}

func TestMangadexGetInfoJapanese(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "jp")
	defer cleanup()

	title, chapter, err := md.GetInfo(md.chapterBase + "/chapter-1")
	require.NoError(t, err)
	require.Equal(t, "テスト", title)
	require.Equal(t, "1", chapter)
}

func TestMangadexGetInfoNoCountry(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "")
	defer cleanup()

	title, chapter, err := md.GetInfo(md.chapterBase + "/chapter-1")
	require.NoError(t, err)
	require.Equal(t, "テスト", title)
	require.Equal(t, "1", chapter)
}

func TestMangadexInitializeMetadataTagsGenres(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "en")
	defer cleanup()

	comic := &core.ComicIssue{
		Source: &core.ComicSource{Name: "test-source", URL: md.chapterBase + "/chapter-1"},
	}
	err := md.Initialize(comic)
	require.NoError(t, err)

	// ensure metadata exists and contains expected genres and tags
	require.NotNil(t, comic.SeriesMetadata)
	require.Contains(t, comic.SeriesMetadata.Genres, "Romance")
	// publicationDemographic should be added to Tags
	require.Contains(t, comic.SeriesMetadata.Tags, "shounen")
	// format tag should include Doujinshi
	require.Contains(t, comic.SeriesMetadata.Tags, "Doujinshi")
	// theme tags like "School Life" are not classified to Tags/Genres by current logic
	require.Contains(t, comic.SeriesMetadata.Tags, "School Life")
}

func TestMangadexInitializeMetadataRating(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "en")
	defer cleanup()

	comic := &core.ComicIssue{
		Source: &core.ComicSource{Name: "test-source", URL: md.chapterBase + "/chapter-1"},
	}
	err := md.Initialize(comic)
	require.NoError(t, err)

	require.NotNil(t, comic.SeriesMetadata)
	require.NotNil(t, comic.SeriesMetadata.CommunityRating)
	require.Equal(t, 4.5, *comic.SeriesMetadata.CommunityRating)
}

func TestMangadexGetAuthorInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/author/author-1"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"result":"ok","data":{"attributes":{"name":"John Doe"}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
	)

	opts := &config.Options{
		URL:            server.URL,
		Country:        "",
		SourceName:     "mangadex.org",
		Logger:         logger.NewLogger(false, nil),
		Client:         client,
		RequestTimeout: config.DefaulltRequestTimeout,
	}

	md := NewMangadex(opts)
	md.apiBase = server.URL

	name, err := md.getAuthorInfo("author-1")
	require.NoError(t, err)
	require.Equal(t, "John Doe", name)
}
