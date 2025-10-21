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
	readAllIssuePath = "/sandman-v2-075-1989/"
	readAllIssueAlt  = "/sandman-v2-_the_deluxe_edition-5-part-6-1989/"
	readAllCategory  = "/category/sandman/"
)

func newReadAllComicsServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		switch r.URL.Path {
		case readAllIssuePath, readAllIssueAlt:
			html := `
                <html>
                    <body>
                        <img src="https://cdn.example.com/sandman/001.jpg" />
                        <img src="https://cdn.example.com/sandman/002.jpg" />
                        <select id="selectbox">
                            <option value="` + base + readAllIssuePath + `"></option>
                            <option value="` + base + readAllIssueAlt + `"></option>
                        </select>
                    </body>
                </html>`
			_, _ = w.Write([]byte(html))
		case readAllCategory:
			html := `
                <html>
                    <body>
                        <ul class="list-story">
                            <a href="` + base + readAllIssuePath + `"></a>
                            <a href="` + base + readAllIssueAlt + `"></a>
                        </ul>
                    </body>
                </html>`
			_, _ = w.Write([]byte(html))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestReadAllComicsScraper(t *testing.T) {
	server := newReadAllComicsServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + readAllIssuePath,
		Logger: logger.NewLogger(false, nil),
	}

	scraper := NewReadallcomics(opts)

	comic := &core.Comic{URLSource: server.URL + readAllIssuePath}
	require.NoError(t, scraper.Initialize(comic))
	require.Equal(t, []string{
		"https://cdn.example.com/sandman/001.jpg",
		"https://cdn.example.com/sandman/002.jpg",
	}, comic.Links)

	// Category listing for All
	opts.All = true
	opts.URL = server.URL + readAllCategory
	scraper = NewReadallcomics(opts)
	issues, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		server.URL + readAllIssuePath,
		server.URL + readAllIssueAlt,
	}, issues)

	// Last issue from category
	opts.All = false
	opts.Last = true
	scraper = NewReadallcomics(opts)
	last, err := scraper.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{server.URL + readAllIssueAlt}, last)
}

func TestReadAllComicsGetInfoParsing(t *testing.T) {
	scraper := NewReadallcomics(&config.Options{})

	tests := []struct {
		url           string
		expectedName  string
		expectedIssue string
	}{
		{"https://readallcomics.com/something-is-killing-the-children-000-2024/", "something is killing the children", "000-2024"},
		{"https://readallcomics.com/sandman-v2-075-1989/", "sandman", "v2-075-1989"},
		{"https://readallcomics.com/sandman-v2-_the_deluxe_edition-5-part-6-1989/", "sandman", "v2-_the_deluxe_edition-5-part-6-1989"},
	}

	for _, tc := range tests {
		name, issue := scraper.GetInfo(tc.url)
		require.Equal(t, tc.expectedName, name)
		require.Equal(t, tc.expectedIssue, issue)
	}
}
