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
	"github.com/stretchr/testify/require"
)

const (
	mangaKakalotChapterPath = "/chapter/manga-title/chapter-2"
	mangaKakalotListPath    = "/manga/manga-title"
)

func newMangaKakalotServer() *httptest.Server {
	chapterHTML := `
		<html>
			<body>
				<div class="breadcrumb">
					<p>
						<span itemprop="itemListElement"><a><span>Home</span></a></span>
						<span itemprop="itemListElement"><a><span>Chapter 2 : My Manga</span></a></span>
					</p>
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
				<div class="chapter-list">
					<div class="row">
						<a href="%s` + mangaKakalotChapterPath + `"></a>
					</div>
					<div class="row">
						<a href="%s/chapter/manga-title/chapter-1"></a>
					</div>
				</div>
			</body>
		</html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == mangaKakalotChapterPath:
			_, _ = fmt.Fprint(w, chapterHTML)
		case r.URL.Path == mangaKakalotListPath:
			base := "http://" + r.Host
			_, _ = fmt.Fprintf(w, listHTMLTemplate, base, base)
		case strings.HasPrefix(r.URL.Path, "/chapter/manga-title/chapter-1"):
			_, _ = fmt.Fprint(w, "<html></html>")
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestMangaKakalotScraper(t *testing.T) {
	server := newMangaKakalotServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + mangaKakalotListPath,
		Source: "mangakakalot.com",
		Logger: logger.NewLogger(false, nil),
	}
	scraper := NewMangaKakalot(opts)

	title, issue := scraper.GetInfo(server.URL + mangaKakalotChapterPath)
	require.Equal(t, "My Manga", title)
	require.Equal(t, "2", issue)

	comic := &core.Comic{URLSource: server.URL + mangaKakalotChapterPath}
	require.NoError(t, scraper.Initialize(comic))
	require.Equal(t, []string{
		"https://cdn.example.com/manga-title/001.jpg",
		"https://cdn.example.com/manga-title/002.jpg",
	}, comic.Links)

	links, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		server.URL + mangaKakalotChapterPath,
		server.URL + "/chapter/manga-title/chapter-1",
	}, links)
}
