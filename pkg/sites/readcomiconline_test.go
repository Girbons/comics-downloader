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
	rcoIssuePath = "/Comic/My-Comic/Issue-2"
	rcoListPath  = "/Comic/My-Comic"
)

func newReadComicOnlineServer() *httptest.Server {
	issueHTML := `
        <html>
            <body>
                <script>push('https://2.bp.blogspot.com/abc123=s1600?')</script>
                <script>push('https://2.bp.blogspot.com/def456=s1600?')</script>
            </body>
        </html>`

	listHTML := `
        <html>
            <body>
                <a href="` + rcoIssuePath + `?id=2">Issue 2</a>
                <a href="/Comic/My-Comic/Issue-1?id=1">Issue 1</a>
            </body>
        </html>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case rcoIssuePath:
			_, _ = fmt.Fprint(w, issueHTML)
		case rcoListPath:
			_, _ = fmt.Fprint(w, listHTML)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestReadComicOnlineScraper(t *testing.T) {
	server := newReadComicOnlineServer()
	defer server.Close()

	originalBase := baseUrl
	baseUrl = server.URL
	defer func() { baseUrl = originalBase }()

	opts := &config.Options{
		URL:    server.URL + rcoIssuePath,
		Logger: logger.NewLogger(false, nil),
	}

	scraper := NewReadComiconline(opts)

	comic := &core.Comic{URLSource: server.URL + rcoIssuePath}
	require.NoError(t, scraper.Initialize(comic))
	require.Equal(t, []string{
		"https://2.bp.blogspot.com/abc123=s1600?",
		"https://2.bp.blogspot.com/def456=s1600?",
	}, comic.Links)

	opts.All = true
	opts.URL = server.URL + rcoListPath
	scraper = NewReadComiconline(opts)
	issues, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		server.URL + "/Comic/My-Comic/Issue-2",
		server.URL + "/Comic/My-Comic/Issue-1",
	}, issues)

	opts.All = false
	opts.Last = true
	scraper = NewReadComiconline(opts)
	last, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{server.URL + "/Comic/My-Comic/Issue-2"}, last)
}
