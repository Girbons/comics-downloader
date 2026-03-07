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
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
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
	md.uploadsBase = server.URL + "/data"

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
		md.uploadsBase + "/HASH/001.png",
		md.uploadsBase + "/HASH/002.png",
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
		md.uploadsBase + "/HASH/001.png",
		md.uploadsBase + "/HASH/002.png",
	}, comic.ImageLinks)
}

func TestMangadexGetInfoEnglish(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "en")
	defer cleanup()

	title, chapter, err := md.GetInfo(md.chapterBase + "/chapter-1")
	require.NoError(t, err)
	require.Equal(t, "Test Manga", title)
	require.Equal(t, "Vol 1 Chapter 1, Start", chapter)
}

func TestMangadexGetInfoJapanese(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "jp")
	defer cleanup()

	title, chapter, err := md.GetInfo(md.chapterBase + "/chapter-1")
	require.NoError(t, err)
	require.Equal(t, "テスト", title)
	require.Equal(t, "Vol 1 Chapter 1, Start", chapter)
}

func TestMangadexGetInfoNoCountry(t *testing.T) {
	md, cleanup := newTestMangadex(t, "series-1", "")
	defer cleanup()

	title, chapter, err := md.GetInfo(md.chapterBase + "/chapter-1")
	require.NoError(t, err)
	require.Equal(t, "テスト", title)
	require.Equal(t, "Vol 1 Chapter 1, Start", chapter)
}
