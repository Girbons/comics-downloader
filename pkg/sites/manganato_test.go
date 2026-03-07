package sites

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/require"
)

const (
	manganatoChapterPath = "/chapter/manga-title/chapter-2"
	manganatoListPath    = "/read/manga-title"
)

func newManganatoServer() *httptest.Server {
	chapterHTML := `
        <html>
            <body>
                <div class="panel-breadcrumb">
                    <a class="a-h">Home</a>
                    <a class="a-h">Chapter 2 : My Manga</a>
                </div>
                <div class="container-chapter-reader">
                    <img src="https://cdn.example.com/manga-title/001.jpg"/>
                    <img src="https://cdn.example.com/manga-title/002.jpg"/>
                </div>
            </body>
        </html>`

	listHTMLTemplate := `
        <html>
            <body>
                <div class="panel-story-chapter-list">
                    <li class="a-h">
                        <a href="%s` + manganatoChapterPath + `"></a>
                    </li>
                    <li class="a-h">
                        <a href="%s/chapter/manga-title/chapter-1"></a>
                    </li>
                </div>
            </body>
        </html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case manganatoChapterPath:
			_, _ = fmt.Fprint(w, chapterHTML)
		case manganatoListPath:
			base := "http://" + r.Host
			_, _ = fmt.Fprintf(w, listHTMLTemplate, base, base)
		case "/chapter/manga-title/chapter-1":
			_, _ = fmt.Fprint(w, "<html></html>")
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestManganatoScraper(t *testing.T) {
	server := newManganatoServer()
	defer server.Close()

	opts := &config.Options{
		URL:            server.URL + manganatoListPath,
		SourceName:     "manganato.com",
		Logger:         logger.NewLogger(false, nil),
		RequestTimeout: config.DefaulltRequestTimeout,
	}

	scraper := NewManganato(opts)

	title, issue, err := scraper.GetInfo(server.URL + manganatoChapterPath)
	require.NoError(t, err)
	require.Equal(t, "My Manga", title)
	require.Equal(t, "2", issue)

	comic := &core.ComicIssue{
		Source: &core.ComicSource{Name: "test-source", URL: server.URL + manganatoChapterPath},
	}
	require.NoError(t, scraper.Initialize(comic))
	require.Equal(t, []string{
		"https://cdn.example.com/manga-title/001.jpg",
		"https://cdn.example.com/manga-title/002.jpg",
	}, comic.ImageLinks)

	links, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		server.URL + manganatoChapterPath,
		server.URL + "/chapter/manga-title/chapter-1",
	}, links)
}
