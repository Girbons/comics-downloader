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
	mangatownIssuePath = "/manga/naruto/v63/c684/"
)

func newMangatownServer() *httptest.Server {
	firstPage := `
        <html>
            <body>
                <div class="page_select">
                    <select>
                        <option>Featured</option>
                        <option>1</option>
                        <option>2</option>
                    </select>
                </div>
                <div id="viewer"><a><img src="//cdn.example.com/naruto/001.jpg"/></a></div>
            </body>
        </html>`

	secondPage := `
        <html>
            <body>
                <div id="viewer"><a><img src="//cdn.example.com/naruto/002.jpg"/></a></div>
            </body>
        </html>`

	listHTML := `
        <html>
            <body>
                <ul class="chapter_list">
                    <a href="/manga/naruto/v63/c684/"></a>
                    <a href="/manga/naruto/v63/c685/"></a>
                </ul>
            </body>
        </html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case mangatownIssuePath:
			_, _ = fmt.Fprint(w, firstPage)
		case "/manga/naruto/v63/c684/1.html":
			_, _ = fmt.Fprint(w, firstPage)
		case "/manga/naruto/v63/c684/2.html":
			_, _ = fmt.Fprint(w, secondPage)
		case "/manga/naruto/v63":
			_, _ = fmt.Fprint(w, listHTML)
		case "/manga/naruto":
			_, _ = fmt.Fprint(w, listHTML)
		case "/manga/naruto/":
			_, _ = fmt.Fprint(w, listHTML)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestMangatownScraper(t *testing.T) {
	server := newMangatownServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + mangatownIssuePath,
		Logger: logger.NewLogger(false, nil),
	}

	scraper := NewMangatown(opts)

	comic := &core.Comic{URLSource: server.URL + mangatownIssuePath}
	require.NoError(t, scraper.Initialize(comic))
	require.Equal(t, []string{
		"https://cdn.example.com/naruto/001.jpg",
		"https://cdn.example.com/naruto/002.jpg",
	}, comic.Links)

	opts.All = true
	scraper = NewMangatown(opts)
	issues, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://mangatown.com/manga/naruto/v63/c684/",
		"https://mangatown.com/manga/naruto/v63/c685/",
	}, issues)

	opts.All = false
	opts.Last = true
	opts.URL = server.URL + "/manga/naruto/"
	scraper = NewMangatown(opts)
	last, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{"https://www.mangatown.com/manga/naruto/v63/c684/"}, last)
}
