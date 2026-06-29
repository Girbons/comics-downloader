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
						"title":{"en":"Test Manga","jp":"テスト"}
					}
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/chapter/chapter-1"):
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{
				"result":"ok",
				"data":{
					"attributes":{"volume":"1","chapter":"1","title":"Start"},
					"relationships":[{"id":"series-1","type":"manga"}]
				}
			}`)
		case strings.HasPrefix(r.URL.Path, "/at-home/server/chapter-1"):
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

func newTestMangadex(t *testing.T) (*Mangadex, func()) {
	t.Helper()

	server := setupMangadexServer()

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
	)

	opts := &config.Options{
		URL:     server.URL + "/title/series-1/naruto",
		Country: "en",
		Source:  "mangadex.org",
		Logger:  logger.NewLogger(false, nil),
		Client:  client,
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
	md, cleanup := newTestMangadex(t)
	defer cleanup()

	md.options.All = true

	links, err := md.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{md.chapterBase + "/chapter-1"}, links)
}

func TestMangadexInitialize(t *testing.T) {
	md, cleanup := newTestMangadex(t)
	defer cleanup()

	comic := &core.Comic{URLSource: md.chapterBase + "/chapter-1"}
	err := md.Initialize(comic)
	require.NoError(t, err)
	require.Equal(t, []string{
		md.uploadsBase + "/HASH/001.png",
		md.uploadsBase + "/HASH/002.png",
	}, comic.Links)
}

func TestMangadexGetInfo(t *testing.T) {
	md, cleanup := newTestMangadex(t)
	defer cleanup()

	title, chapter := md.GetInfo(md.chapterBase + "/chapter-1")
	require.Equal(t, "Test Manga", title)
	require.Equal(t, "Vol 1 Chapter 1, Start", chapter)
}
