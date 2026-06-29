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
	comicExtraIssueFullPath = "/batman-unseen/issue-5/full"
	comicExtraIssueSimple   = "/batman-unseen/issue-5"
	comicExtraLastIssuePath = "/batman-unseen/issue-4/full"
	comicExtraListPath      = "/comic/batman-unseen"
)

func newComicExtraServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "http://" + r.Host
		switch r.URL.Path {
		case comicExtraIssueFullPath, comicExtraIssueSimple, "/batman-unseen/issue-4":
			html := `
                <html>
                    <body>
                        <img src="https:\/\/cdn.example.com\/batman%3Fpage%3D1" />
                        <img src="https://cdn.example.com/batman%3Fpage%3D2" />
                    </body>
                </html>`
			_, _ = w.Write([]byte(html))
		case comicExtraLastIssuePath:
			html := `
                <html>
                    <body>
                        <select>
                            <option value="` + base + comicExtraIssueFullPath + `"></option>
                        </select>
                    </body>
                </html>`
			_, _ = w.Write([]byte(html))
		case comicExtraListPath:
			html := `
                <html>
                    <body>
                        <div class="episode-list">
                            <a href="` + comicExtraIssueFullPath + `"></a>
                            <a href="/batman-unseen/issue-4"></a>
                        </div>
                    </body>
                </html>`
			_, _ = w.Write([]byte(html))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestComicExtraScraper(t *testing.T) {
	server := newComicExtraServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + comicExtraIssueFullPath,
		Logger: logger.NewLogger(false, nil),
	}

	comicextra := NewComicextra(opts)

	comic := &core.Comic{URLSource: server.URL + comicExtraIssueFullPath}
	require.NoError(t, comicextra.Initialize(comic))
	require.Equal(t, []string{
		"https://cdn.example.com/batman?page=1",
		"https://cdn.example.com/batman?page=2",
	}, comic.Links)
}

func TestComicExtraRetrieveIssueLinksAll(t *testing.T) {
	server := newComicExtraServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + comicExtraListPath,
		All:    true,
		Logger: logger.NewLogger(false, nil),
	}

	comicextra := NewComicextra(opts)
	issues, err := comicextra.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{
		comicExtraIssueFullPath,
		"/batman-unseen/issue-4/full",
	}, issues)
}

func TestComicExtraRetrieveLastIssue(t *testing.T) {
	server := newComicExtraServer()
	defer server.Close()

	opts := &config.Options{
		URL:    server.URL + comicExtraLastIssuePath,
		Last:   true,
		Logger: logger.NewLogger(false, nil),
	}

	comicextra := NewComicextra(opts)
	issues, err := comicextra.RetrieveIssueLinks()
	require.NoError(t, err)
	require.Equal(t, []string{server.URL + comicExtraIssueFullPath}, issues)
}
