package sites

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Girbons/comics-downloader/internal/logger"
	"github.com/Girbons/comics-downloader/pkg/config"
	"github.com/Girbons/comics-downloader/pkg/core"
	"github.com/stretchr/testify/require"
)

const (
	mangareaderIssuePath = "/naruto/1/"
	mangareaderBasePath  = "/naruto"
	mangareaderLastIssue = "/naruto/700"
)

func newMangareaderServer() *httptest.Server {
	issueHTML := `
        <html>
            <body>
                <img data-src="https://cdn.example.com/naruto/001.jpg" />
                <img data-src="https://cdn.example.com/naruto/002.jpg" />
            </body>
        </html>`

	baseHTML := `
        <html>
            <body>
                <table class="d48">
                    <tr><td><a href="/naruto/1/">Chapter 1</a></td></tr>
                    <tr><td><a href="/naruto/2/">Chapter 2</a></td></tr>
                </table>
                <ul class="d44">
                    <li><a href="` + mangareaderLastIssue + `">Latest</a></li>
                </ul>
            </body>
        </html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case mangareaderIssuePath:
			_, _ = w.Write([]byte(issueHTML))
		case "/naruto/2/":
			_, _ = w.Write([]byte(issueHTML))
		case mangareaderBasePath:
			_, _ = w.Write([]byte(baseHTML))
		case mangareaderLastIssue:
			_, _ = w.Write([]byte("<html></html>"))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestMangareaderScraper(t *testing.T) {
	server := newMangareaderServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + mangareaderIssuePath,
		Logger: logger.NewLogger(false, nil),
	}

	scraper := NewMangareader(opts)

	comic := &core.Comic{URLSource: server.URL + mangareaderIssuePath}
	require.NoError(t, scraper.Initialize(comic))
	require.Equal(t, []string{
		"https://cdn.example.com/naruto/001.jpg",
		"https://cdn.example.com/naruto/002.jpg",
	}, comic.Links)

	opts.All = true
	opts.URL = server.URL + mangareaderIssuePath
	scraper = NewMangareader(opts)
	issues, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://mangareader.tv" + mangareaderIssuePath,
		"https://mangareader.tv/naruto/2/",
	}, issues)

	opts.Last = true
	opts.All = false
	opts.URL = server.URL + mangareaderBasePath
	scraper = NewMangareader(opts)
	lastIssues, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{"https://mangareader.tv" + mangareaderLastIssue}, lastIssues)
}
